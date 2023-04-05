package controllers

import (
	"Footware-Ecommerce/models"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't find product")
	ErrUserIDIsNotValid   = errors.New("user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add product to cart")
	ErrCantRemoveItem     = errors.New("cannot remove item from cart")
	ErrCantGetItem        = errors.New("cannot get item from cart ")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductStockByAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Query("productID")
		if productID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}
		product, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var sizedetailsStocks models.SizeDetails
		if err = c.BindJSON(&sizedetailsStocks); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: product}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$size"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$product_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := ProductCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var sizedetailsinfo []bson.M
		if err = pointcursor.All(ctx, &sizedetailsinfo); err != nil {
			panic(err)
		}
		var sizedetails int32
		for _, address_no := range sizedetailsinfo {
			count := address_no["count"]
			sizedetails = count.(int32)
		}
		if sizedetails < 1 {
			filter := bson.D{primitive.E{Key: "_id", Value: product}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "size", Value: sizedetailsStocks}}}}
			_, err := ProductCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not Allowed ")
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(201, "Stock Successfully Added.")
	}

}

func UpdateProductStock() gin.HandlerFunc {
	return func(c *gin.Context) {
		product := c.Query("productID")
		if product == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid"})
			c.Abort()
			return
		}
		product_id, err := primitive.ObjectIDFromHex(product)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var updateProductStock models.SizeDetails
		if err := c.BindJSON(&updateProductStock); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: product_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "size.0.UK5_5", Value: updateProductStock.UK55},
			{Key: "size.0.UK6", Value: updateProductStock.UK6}, {Key: "size.0.UK6_5", Value: updateProductStock.UK65},
			{Key: "size.0.UK7", Value: updateProductStock.UK7}, {Key: "size.0.UK7_5", Value: updateProductStock.UK75},
			{Key: "size.0.UK8", Value: updateProductStock.UK8}, {Key: "size.0.UK8_5", Value: updateProductStock.UK85},
			{Key: "size.0.UK9", Value: updateProductStock.UK9}, {Key: "size.0.UK9_5", Value: updateProductStock.UK95},
			{Key: "size.0.UK10", Value: updateProductStock.UK10}, {Key: "size.0.UK10_5", Value: updateProductStock.UK105},
			{Key: "size.0.UK11", Value: updateProductStock.UK11}, {Key: "size.0.UK11_5", Value: updateProductStock.UK115},
			{Key: "size.0.UK12", Value: updateProductStock.UK12}}}}
		_, err = ProductCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Panic(err)
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Updated the Product Stock")
	}
}
