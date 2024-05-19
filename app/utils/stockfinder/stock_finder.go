package stockfinder

import (
	"errors"
	"fmt"
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

func DownloadStocksData(database *gorm.DB, rangeMap map[string]string) error {
	if os.Getenv("DOWNLOAD_STOCKS_DATA") == "false" {
		return errors.New("DOWNLOAD_STOCKS_DATA env. variable is false")
	}

	stocks, err := db.GetAllStocks(database)
	if err != nil {
		return err
	}

	for _, stock := range stocks {
		for interval, period1 := range rangeMap {
			fmt.Printf("%s Downloading (interval: %s)\n", stock.Code, interval)
			err = stock.DownloadYahooCSV(interval, period1, "1716123973")
			if err != nil {
				fmt.Printf("Error: Downloading Yahoo CSV (%s): %s\n", err.Error(), stock.Code)
				continue
			}
		}
	}

	return nil
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
