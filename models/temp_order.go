package models

import "time"

type TempOrder struct {
	COrderID uint      `json:"c_order_id" gorm:"column:c_order_id"`
	Created  time.Time `json:"created" gorm:"column:created"`
}
