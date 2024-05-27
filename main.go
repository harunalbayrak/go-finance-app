package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/harunalbayrak/go-finance-app/app/chart"
	"github.com/harunalbayrak/go-finance-app/app/configs"
	"github.com/harunalbayrak/go-finance-app/app/telegram"
	"github.com/harunalbayrak/go-finance-app/app/utils/stockfinder"
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() error {
	err := godotenv.Load(".env")

	return err
}

func StartApp() {
	engine := html.New("./views", ".html")
	config := configs.FiberConfig(engine)
	app := fiber.New(config)
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})
	log.Fatal(app.Listen(":3000"))
}

func main() {
	err := LoadEnvironmentVariables()
	if err != nil {
		log.Fatal("env loading error", err)
	}

	// database, err := db.CreateDB()
	// if err != nil {
	// 	log.Fatal("creating db error", err)
	// }

	// err = stockfinder.CreateStocksOnDatabase(database)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	rangeMap := make(map[string]string)
	rangeMap["1d"] = "643122373"

	// err = stockfinder.DownloadStocksData(database, rangeMap)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	err = stockfinder.DownloadStockData("EKIZ", rangeMap)
	if err != nil {
		fmt.Println(err.Error())
	}

	chart.StockChart("EKIZ", 90, 1)

	fmt.Println(time.Now().Unix())

	telegram.SendMessage("ASDASD")
	telegram.SendPhoto("charts/EKIZ.png")
}

// func other() {
// StartApp()
// stocks, _ := db.GetAllStocks(database)

// cookie, err := yahoo.GetCookie()
// crumb, err := yahoo.GetCrumb(cookie)

// start := time.Now()
// for i, stock := range stocks {
// 	if i == 1 {
// 		break
// 	}

// 	fmt.Println("stock:", stock)
// 	stock := models.Stock{Code: "EKIZ"}
// 	// yahooChart, _ := stock.GetYahooChart("1d", "5d")
// 	// fmt.Println("Stock:", yahooChart.Chart.Result[0].Meta.Symbol, "=", yahooChart.Chart.Result[0].Meta.RegularMarketPrice)

// 	// yahooQuoteResponse, _ := stock.GetYahooQuoteResponse(cookie, crumb)
// 	// fmt.Printf("QuoteResponse: %+v\n", yahooQuoteResponse.QuoteResponse.Result)

// 	stock.DownloadYahooCSV("1d", "643122373", "1716123973")

// }
// timeElapsed := time.Since(start)
// fmt.Printf("The `for` loop took %s\n", timeElapsed)

// intervalStr := os.Getenv("REFRESH_INTERVAL")
// interval, err := strconv.Atoi(intervalStr)
// if err != nil {
// 	log.Fatal("interval error", err)
// }
// if interval < 1 || interval > 300 {
// 	log.Fatal("interval size error", err)
// }
// go scheduler.Run(time.Duration(interval)*time.Minute, time.Minute)

// // telegram.SendMessage(string(t.Render()))

// StartApp()
// }
