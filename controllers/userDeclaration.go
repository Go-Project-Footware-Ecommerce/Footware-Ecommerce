package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"Footware-Ecommerce/database"
	helper "Footware-Ecommerce/helpers"
	"Footware-Ecommerce/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")          //To store data in Users collection in mongoDB database
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products") // TO store data in Product collection in mongoDB database
var Validate = validator.New()

// GenerateFromPassword returns the bcrypt hash of the password at the given cost.
// If the cost given is less than MinCost, the cost will be set to DefaultCost, instead.
// Use CompareHashAndPassword, as defined in this package, to compare the returned hashed password
// with its cleartext version. GenerateFromPassword does not accept passwords longer than 72 bytes,
// which is the longest password bcrypt will operate on.
func HashPassword(password string) string { // function used for convert the string passward in bcrypt hash
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// verify the password enter by user.user enter in string format then we convert into bcryt hash and compare it into our databse.
func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

// signup of user or admin
func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		//context is a standard package of Golang that makes it easy to pass request-scoped values, cancelation signals,
		//and deadlines across API boundaries to all the goroutines involved in handling a request.
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) //to running backgroup process.when deal with the database.
		c.HTML(http.StatusOK, "signup.html", nil)
		var user models.User //custom datatype

		if err := c.BindJSON(&user); err != nil { //respone received from postman.
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user) // validate the datatype requirement.
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		//increase the count when email or phone no exit in database or already user signup.
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		//convert the password by using bcrypt adaptive hashing algorithm
		password := HashPassword(user.Password)
		user.Password = password

		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone}) //bson.M == An unordered representation of a BSON document (map)
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
		}
		//Internal Server Error
		//The HyperText Transfer Protocol (HTTP) 500 Internal Server Error server error response code indicates that
		//the server encountered an unexpected condition that prevented it from fulfilling the request.
		//This error response is a generic "catch-all" response.

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		}

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		//Primitive Data Types: A primitive data type is pre-defined by the programming language.
		//The size and type of variable values are specified, and it has no additional methods.
		//What is the use of ObjectId?
		//MongoDB uses ObjectIds as the default value of _id field of each document, which is generated during the
		//creation of any document. Object ID is treated as the primary key within any MongoDB collection.
		//It is a unique identifier for each document or record.
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.First_Name, user.Last_Name, user.User_type, user.User_ID)
		user.Token = token //generate the token from given data by using bcrypt package
		user.Refresh_Token = refreshToken
		//Refresh token: The refresh token is used to generate a new access token.
		//Typically, if the access token has an expiration date, once it expires, the user would have to
		//authenticate again to obtain an access token
		resultInsertionNumber, insertErr := UserCollection.InsertOne(ctx, user) //insert data into database
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel() //closed the database.
		//A Context is a standard Go data value that can report whether the overall operation it represents
		//has been canceled and is no longer needed
		c.JSON(http.StatusCreated, resultInsertionNumber)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser) //find email in database
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password) //check user entered password is correct or not.
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email not found.Please enter correct email OR please SIGNUP"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(foundUser.Email, foundUser.First_Name, foundUser.Last_Name, foundUser.User_type, foundUser.User_ID)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
		err = UserCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_ID}).Decode(&foundUser) // find user details in database.

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func ProductAddByAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product //customer product details datatype
		defer cancel()
		if err := c.BindJSON(&products); err != nil { //received data from postman in format of json and convert in go datatype means unmarshalling.
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		products.Size = make([]models.SizeDetails, 0)
		products.Review = make([]models.UserReview, 0)
		_, anyerr := ProductCollection.InsertOne(ctx, products) //insert data into the database.
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully added our Product!!")
	}
}

func AllProductList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.D{{}}) //gives all product data from database
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
			return
		}
		err = cursor.All(ctx, &productlist) //data recived in the form of json we need to comvert into go structure.
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			// Don't forget to log errors. I log them really simple here just
			// to get the point across.
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, productlist)

	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchproducts []models.Product
		queryParam := c.Query("name") //name of product which you want to seach
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}}) //search in database
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchproducts) //search result received and convert into go datatype
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx) //close the database query.
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchproducts)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		//To check user type
		//if user_type is admin it will gives you all users data from database.
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage")) //view only 10 user in one page
		if err != nil || recordPerPage < 1 {
			recordPerPage = 30
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := UserCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := UserCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user) //find user_id in database
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// func HomePage () gin.HandlerFunc{
// 	return func(c *gin.Context) {
// 		c.HTML(http.StatusOK,"homepage.html",nil)
// 	}
// }
