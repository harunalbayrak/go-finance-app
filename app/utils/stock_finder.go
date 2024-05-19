package utils

import (
	"errors"
	"log"
	"os"

	"github.com/gocolly/colly"
	"github.com/harunalbayrak/go-finance-app/app/db"
	"github.com/harunalbayrak/go-finance-app/app/models"
	"gorm.io/gorm"
)

func CreateStocksOnDatabase(database *gorm.DB) error {
	if os.Getenv("FIND_STOCKS") == "false" {
		return errors.New("FIND_STOCKS env. variable is false")
	}

	// AutoMigrate for player table
	err := database.AutoMigrate(&models.Stock{})
	if err != nil {
		log.Fatal("migrating db error", err)
	}

	stocks, err := FindAllStocks()
	if err != nil {
		log.Fatal("finding stocks error", err)
	}

	db.CreateStocks(database, stocks)

	return err
}

func FindAllStocks() ([]models.Stock, error) {
	stocks := make([]models.Stock, 0)
	urlBase := os.Getenv("STOCK_FINDER_URL_BASE")
	urlPrefix := os.Getenv("STOCK_FINDER_URL_PREFIX")

	c := colly.NewCollector(
		colly.AllowedDomains(urlBase),
	)

	c.OnHTML("tr", func(e *colly.HTMLElement) {
		code := e.ChildText("b")
		if code != "" {
			stocks = append(stocks, models.Stock{Code: code})
		}
	})

	c.Visit("https://" + urlBase + urlPrefix)

	return stocks, nil
}
