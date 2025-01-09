package controllers

import (
	"bima-go/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderDetail struct {
	TableName              string
	ID                     int
	COrderlineID           int
	COrderID               int
	CBPartnerID            int
	ItemID                 int
	FactoryID              int
	SoID                   int
	CategoryID             int
	PurchaseNumber         string
	SupplierName           string
	Category               string
	ItemCode               string
	UOM                    string
	QtyAllocation          float64
	POBuyer                string
	LCDate                 string
	IsRecycle              string
	StatusLC               string
	StdPrecision           int
	SoOrderType            string
	Season                 string
	WarehousePlace         string
	SoDocTypeID            int
	PromiseDate            time.Time
	JobOrder               string
	PurchaseDocumentTypeID int
	IsFabric               string
	Color                  string
	UOMID                  int
	ItemDesc               string
	LastUpdatePO           time.Time
	IsMrpExists            bool
}

func GetAutoPr(c *gin.Context) {
	lcDate := time.Now().AddDate(0, -2, 0)
	lastUpdatePo := time.Now().AddDate(0, -2, 0)

	rows, err := models.DB.Table("bw_wms_auto_pr as bwap").
		Select("*").
		Where("bwap.lc_date >= ?", lcDate).
		Where("bwap.last_update_po >= ?", lastUpdatePo).
		Where("bwap.is_mrp_exists = ?", true).
		Rows()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	var orderDetails []OrderDetail
	for rows.Next() {
		var orderDetail OrderDetail
		models.DB.ScanRows(rows, &orderDetail)
		orderDetails = append(orderDetails, orderDetail)
		if len(orderDetails) >= 100 {
			break
		}
	}

	c.JSON(http.StatusOK, orderDetails)
}
