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
	}

	return router
}
