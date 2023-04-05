package main

import (
	"log"

	"Footware-Ecommerce/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8000"

	router := gin.New()              //To create a new router, use the New() function
	router.Use(gin.Logger())         //Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
	routes.EveyoneShowRoutes(router) //Everyone can see this APIs
	routes.AdminAuthRoutes(router)   //Required ADMIN token to use this APIs
	routes.UserAuthRoutes(router)    //Required USER token to use this APIs

	log.Fatal(router.Run(":" + port))
}
