package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TempOrder struct {
	COrderID int       `gorm:"column:c_order_id"`
	Created  time.Time `gorm:"column:created"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	bimaOrderTempQuery := db.Table("bima_order_temp").
		Select("c_order_id, DATE(created) AS created").
		Where("issotrx = ?", "N").
		Where("created >= ?", "2024-01-01")

	bimaOrderLineTempQuery := db.Table("bima_orderline_temp").
		Select("c_order_id, DATE(created) AS created").
		Where("created >= ?", "2024-01-01")

	var tempOrders []TempOrder

	err = db.Raw("? UNION ?", bimaOrderTempQuery, bimaOrderLineTempQuery).Scan(&tempOrders).Error
	if err != nil {
		log.Fatal(err)
	}
	for _, order := range tempOrders {
		fmt.Printf("COrderID: %d, Created: %s\n", order.COrderID, order.Created.Format("2006-01-02"))
	}
}
