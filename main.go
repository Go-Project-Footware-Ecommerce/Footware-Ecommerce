package main

import (
	"log"

	"Footware-Ecommerce/routes"

	"github.com/gin-gonic/gin"
)

/*
// variable for template
var tmpl *template.Template

// init function to initialize this object of template
func init() {
	//specify the directory that will have my all html code
	//this go will parse this folder and * means all the html files under this folder
	tmpl = template.Must(template.ParseGlob("template/*.html"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//here i will execute my index file that is my home file
	tmpl.ExecuteTemplate(w, "index.html", nil)
}
*/

func main() {
	port := "8000"

	router := gin.New()      //To create a new router, use the New() function
	router.Use(gin.Logger()) //Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
	//router.LoadHTMLGlob("template/*.html")
	routes.EveyoneShowRoutes(router) //Everyone can see this APIs
	routes.AdminAuthRoutes(router)   //Required ADMIN token to use this APIs
	routes.UserAuthRoutes(router)    //Required USER token to use this APIs

	log.Fatal(router.Run(":" + port))
}
