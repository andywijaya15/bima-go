package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestSpeed))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/orders", func(c *gin.Context) {
		bimaOrderTempQuery := db.Table("bima_order_temp").
			Select("c_order_id, DATE(created) AS created").
			Where("issotrx = ?", "N").
			Where("created::date >= ?", "2024-01-01")

		// bimaOrderLineTempQuery := db.Table("bima_orderline_temp").
		// 	Select("c_order_id, DATE(created) AS created").
		// 	Where("created::date >= ?", "2024-01-01")

		var tempOrders []TempOrder

		// err = db.Raw("SELECT DISTINCT * FROM (? UNION ALL ?) tbl", bimaOrderTempQuery, bimaOrderLineTempQuery).Scan(&tempOrders).Error
		err = db.Raw("SELECT DISTINCT * FROM (?) tbl", bimaOrderTempQuery).Scan(&tempOrders).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"status":  http.StatusOK,
			"message": "List Data",
			"data":    tempOrders,
			"count":   len(tempOrders),
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
