package controllers

import (
	"bima-go/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func deleteOrder(tx *gorm.DB, tableName, cOrderId string) error {
	err := tx.Table(tableName).
		Where("issotrx = ?", "N").
		Where("c_order_id = ?", cOrderId).
		Delete(nil).Error

	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", tableName, err)
	}
	return nil
}

func DeletePurchaseOrder(c *gin.Context) {
	cOrderId := c.Query("c_order_id")
	if cOrderId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "c_order_id is required"})
		return
	}
	tx := models.DB.Begin()
	if tx.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction", "details": tx.Error.Error()})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if err := tx.Table("bima_order_temp").
		Where("issotrx = ?", "N").
		Where("c_order_id = ?", cOrderId).
		Delete(nil).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Table("bima_orderline_temp").
		Where("c_order_id = ?", cOrderId).
		Delete(nil).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order deleted successfully"})
}
