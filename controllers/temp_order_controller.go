package controllers

import (
	"bima-go/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetChangedPurchaseOrders(c *gin.Context) {
	var tempOrders []models.TempOrder
	threshold := time.Now()

	// Subtract 3 months from the threshold
	threshold = threshold.AddDate(0, -3, 0)
	bimaOrderTempQuery := models.DB.Table("bima_order_temp").
		Select("c_order_id, DATE(created) AS created").
		Where("issotrx = ?", "N").
		Where("created::date >= ?", threshold)

	bimaOrderLineTempQuery := models.DB.Table("bima_orderline_temp").
		Select("c_order_id, DATE(created) AS created").
		Where("created::date >= ?", threshold)

	err := models.DB.Raw("SELECT DISTINCT * FROM (? UNION ALL ?) tbl", bimaOrderTempQuery, bimaOrderLineTempQuery).Scan(&tempOrders).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tempOrders)
}
