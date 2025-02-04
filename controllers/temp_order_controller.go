package controllers

import (
	"bima-go/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TempOrder struct {
	COrderID  uint      `gorm:"column:c_order_id"`
	Created   time.Time `gorm:"column:created"`
	IsActive  string    `gorm:"column:isactive"`
	DocStatus string    `gorm:"column:docstatus"`
}

// tanpa concurrency
func GetChangedPurchaseOrders(c *gin.Context) {
	var tempOrders []TempOrder
	threshold := time.Now()

	threshold = threshold.AddDate(0, -3, 0)
	bimaOrderTempQuery := models.DB.Table("bima_order_temp as bot").
		Joins("JOIN c_order po ON po.c_order_id = bot.c_order_id").
		Select("bot.c_order_id, DATE(bot.created) AS created, po.isactive, po.docstatus").
		Where("bot.issotrx = ?", "N").
		Where("bot.created::date >= ?", threshold)

	bimaOrderLineTempQuery := models.DB.Table("bima_orderline_temp as bolt").
		Joins("JOIN c_order po ON po.c_order_id = bolt.c_order_id").
		Select("bolt.c_order_id, DATE(bolt.created) AS created, po.isactive, po.docstatus").
		Where("bolt.created::date >= ?", threshold)

	err := models.DB.Raw("SELECT DISTINCT * FROM (? UNION ALL ?) tbl", bimaOrderTempQuery, bimaOrderLineTempQuery).Scan(&tempOrders).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, tempOrders)
}

// dengan concurrency
func GetChangedPurchaseOrdersConcurrency(c *gin.Context) {
	threshold := time.Now().AddDate(0, -3, 0)
	tempOrders := make(map[uint]TempOrder)

	// Channels for goroutines
	bimaOrderTempChan := make(chan []TempOrder, 1)
	bimaOrderLineTempChan := make(chan []TempOrder, 1)
	errChan := make(chan error, 2)

	// Goroutine for the first query
	go func() {
		var orders []TempOrder
		err := models.DB.Table("bima_order_temp as bot").
			Joins("LEFT JOIN c_order po ON po.c_order_id = bot.c_order_id").
			Select("bot.c_order_id, DATE(bot.created) AS created, po.isactive, po.docstatus").
			Where("bot.issotrx = ?", "N").
			Where("bot.created::date >= ?", threshold).
			Scan(&orders).Error
		if err != nil {
			errChan <- err
			return
		}
		bimaOrderTempChan <- orders
	}()

	// Goroutine for the second query
	go func() {
		var orders []TempOrder
		err := models.DB.Table("bima_orderline_temp as bolt").
			Joins("LEFT JOIN c_order po ON po.c_order_id = bolt.c_order_id").
			Select("bolt.c_order_id, DATE(bolt.created) AS created, po.isactive, po.docstatus").
			Where("bolt.created::date >= ?", threshold).
			Scan(&orders).Error
		if err != nil {
			errChan <- err
			return
		}
		bimaOrderLineTempChan <- orders
	}()

	// Collect results
	var (
		bimaOrderTempResult     []TempOrder
		bimaOrderLineTempResult []TempOrder
	)

	for i := 0; i < 2; i++ {
		select {
		case orders := <-bimaOrderTempChan:
			bimaOrderTempResult = orders
		case orders := <-bimaOrderLineTempChan:
			bimaOrderLineTempResult = orders
		case err := <-errChan:
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Merge results and remove duplicates
	for _, order := range append(bimaOrderTempResult, bimaOrderLineTempResult...) {
		tempOrders[order.COrderID] = order
	}

	// Convert map back to slice
	result := make([]TempOrder, 0, len(tempOrders))
	for _, order := range tempOrders {
		result = append(result, order)
	}

	c.JSON(http.StatusOK, result)
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
