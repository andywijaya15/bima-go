package controllers

import (
	"bima-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetChangedPurchaseOrders(c *gin.Context) {
	var tempOrders []models.TempOrder
	bimaOrderTempQuery := models.DB.Table("bima_order_temp").
		Select("c_order_id, DATE(created) AS created").
		Where("issotrx = ?", "N").
		Where("created::date >= ?", "2024-01-01")

	bimaOrderLineTempQuery := models.DB.Table("bima_orderline_temp").
		Select("c_order_id, DATE(created) AS created").
		Where("created::date >= ?", "2024-01-01")

	err := models.DB.Raw("SELECT DISTINCT * FROM (? UNION ALL ?) tbl", bimaOrderTempQuery, bimaOrderLineTempQuery).Scan(&tempOrders).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tempOrders)
}
