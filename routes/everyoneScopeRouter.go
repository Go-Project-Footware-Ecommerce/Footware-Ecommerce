package routes

import (
	"Footware-Ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func EveyoneShowRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users/signup", controllers.Signup())
	incomingRoutes.POST("/users/signup", controllers.Signup())              //Signup API
	incomingRoutes.POST("/users/login", controllers.Login())                //Login API
	incomingRoutes.GET("/productview", controllers.AllProductList())  //To view all product in database
	incomingRoutes.GET("/search", controllers.SearchProductByQuery()) //Search product by using PRODUCT NAME
}
