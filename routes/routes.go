package routes

import (
	"bima-go/controllers"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.BestCompression))

	v1 := router.Group("/v1")
	{
		v1.GET("/get-changed-purchase-orders", controllers.GetChangedPurchaseOrders)
		v1.GET("/get-changed-purchase-orders-concurrency", controllers.GetChangedPurchaseOrdersConcurrency)
		v1.POST("/delete-purchase-order", controllers.DeletePurchaseOrder)

		v1.GET("/get-auto-pr", controllers.GetAutoPr)
	}

	return router
}
