package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"Footware-Ecommerce/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProductReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		product_id := c.Query("productID")
		if product_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}
		product, err := primitive.ObjectIDFromHex(product_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		var addReview models.UserReview
		addReview.UserID = userQueryID
		if err = c.BindJSON(&addReview); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: product}}}}
		// unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$reviews"}}}}
		// group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$product_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		// _, err = ProductCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		// if err != nil {
		// 	c.IndentedJSON(500, "Internal Server Error")
		// }
		// var reviewinfo []bson.M
		// if err = pointcursor.All(ctx, &reviewinfo); err != nil {
		// 	panic(err)
		// }
		filter := bson.D{primitive.E{Key: "_id", Value: product}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "review", Value: addReview}}}}
		_, err = ProductCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			fmt.Println(err)
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(201, "Review Added")
	}

}
