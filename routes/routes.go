package routes

import (
	"bima-go/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/get-changed-purchase-orders", controllers.GetChangedPurchaseOrders)
		v1.POST("/delete-purchase-order", controllers.DeletePurchaseOrder)
	}

	return router
}
