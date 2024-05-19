package db

import (
	"fmt"
	"os"

	"github.com/harunalbayrak/go-finance-app/app/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_CONNECTION") + ")/test?charset=utf8mb4"
	fmt.Println("Dsn:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: log.Default.LogMode(log.Info),
	})

	return db, err
}

func CreateStock(db *gorm.DB, stock *models.Stock) error {
	if err := db.Where("code = ?", stock.Code).Updates(stock).FirstOrCreate(stock).Error; err != nil {
		return nil
	}

	return nil
}

func CreateStocks(db *gorm.DB, stocks []models.Stock) error {
	for _, stock := range stocks {
		CreateStock(db, &stock)
	}

	return nil
}

func GetAllStocks(db *gorm.DB) ([]models.Stock, error) {
	stocks := []models.Stock{}
	db.Find(&stocks)

	return stocks, nil
}
