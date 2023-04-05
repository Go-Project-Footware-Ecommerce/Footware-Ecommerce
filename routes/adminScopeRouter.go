package routes

import (
	controller "Footware-Ecommerce/controllers"
	"Footware-Ecommerce/middleware"

	"github.com/gin-gonic/gin"
)

func AdminAuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/admin/allusers", controller.GetUsers())                       // done
	incomingRoutes.GET("/admin/user/:user_id", controller.GetUser())                   // done
	incomingRoutes.POST("/admin/addproduct", controller.ProductAddByAdmin())           // done
	incomingRoutes.POST("/admin/addproductstock", controller.AddProductStockByAdmin()) // done
	incomingRoutes.PUT("/admin/updateproductstock", controller.UpdateProductStock())   //
}
