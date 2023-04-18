package routes

import (
	"Footware-Ecommerce/controllers"
	"Footware-Ecommerce/database"
	"Footware-Ecommerce/middleware"

	"github.com/gin-gonic/gin"
)

func UserAuthRoutes(incomingRoutes *gin.Engine) {
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/addtocart", app.AddToCart())                      //done
	incomingRoutes.DELETE("/removeitem", app.RemoveItem())                 //done
	incomingRoutes.GET("/listcart", controllers.GetItemFromCart())         //done
	incomingRoutes.POST("/addaddress", controllers.AddAddress())           //done
	incomingRoutes.PUT("/edithomeaddress", controllers.EditHomeAddress())  //done
	incomingRoutes.PUT("/editworkaddress", controllers.EditWorkAddress())  //done
	incomingRoutes.DELETE("/deleteaddresses", controllers.DeleteAddress()) //done
	incomingRoutes.GET("/checkouttocart", app.BuyFromCart())               //done
	incomingRoutes.GET("/buynow", app.InstantBuy())                        //done
	incomingRoutes.POST("/productreview", controllers.ProductReviews())    //done
}
