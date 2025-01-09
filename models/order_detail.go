package models

import "time"

type OrderDetail struct {
	TableName              string
	ID                     int
	COrderlineID           int
	COrderID               int
	CBPartnerID            string
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
	LCDate                 time.Time
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
}
