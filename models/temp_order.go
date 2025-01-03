package models

import "time"

type TempOrder struct {
	COrderID uint      `gorm:"column:c_order_id"`
	Created  time.Time `gorm:"column:created"`
}
