package controllers

import (
	"bima-go/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAutoPr(c *gin.Context) {
	var orderDetails []models.OrderDetail
	lcDate := time.Now().AddDate(0, -3, 0)
	lastUpdatePo := time.Now().AddDate(0, -4, 0)

	err := models.DB.Table("bw_wms_auto_pr as bwap").
		Select("*").
		// Where("bwap.lc_date >= ?", lcDate).
		// Where("bwap.last_update_po >= ?", lastUpdatePo).
		Where("bwap.is_mrp_exists = ?", true).
		Where(models.OrderDetail{
			LCDate:       lcDate,
			LastUpdatePO: lastUpdatePo,
		}).
		Limit(1).
		Scan(&orderDetails).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, orderDetails)
}
