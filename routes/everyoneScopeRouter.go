package routes

import (
	"Footware-Ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func EveyoneShowRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.Signup())              //Signup API
	incomingRoutes.POST("/users/login", controllers.Login())                //Login API
	incomingRoutes.GET("/users/productview", controllers.AllProductList())  //To view all product in database
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery()) //Search product by using PRODUCT NAME
}
