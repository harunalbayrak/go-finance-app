package main

import (
	"fmt"
	"log"
	"time"

	"github.com/harunalbayrak/go-finance-app/app/db"
	"github.com/harunalbayrak/go-finance-app/app/models"
	"github.com/harunalbayrak/go-finance-app/app/utils"
	"github.com/harunalbayrak/go-finance-app/pkg/yahoo"
	"github.com/joho/godotenv"
)

// "github.com/gofiber/fiber/v2"
// "github.com/gofiber/template/django/v3"
// "github.com/harunalbayrak/go-finance/app/db"
// "github.com/harunalbayrak/go-finance/app/models"
// "github.com/harunalbayrak/go-finance/pkg/configs"
// "github.com/harunalbayrak/go-finance/pkg/routes"
// "github.com/harunalbayrak/go-finance/pkg/scheduler"
// "github.com/harunalbayrak/go-finance/pkg/utils"
// "github.com/harunalbayrak/go-finance/pkg/yahoo"
// "github.com/jedib0t/go-pretty/table"

func LoadEnvironmentVariables() error {
	err := godotenv.Load(".env")

	return err
}

// func StartApp() {
// 	engine := django.New("./web/views", ".django")

// 	// AddFunc adds a function to the template's global function map.
// 	engine.AddFunc("greet", func(name string) string {
// 		return "Hello, " + name + "!"
// 	})

// 	config := configs.FiberConfig(engine)

// 	// After you created your engine, you can pass it to Fiber's Views Engine
// 	app := fiber.New(config)

// 	routes.PublicRoutes(app) // Register a public routes for app.

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.Render("index", fiber.Map{
// 			"stocks": []models.Stock{{Code: "asd2"}, {Code: "asd"}},
// 		})
// 	})

// 	routes.NotFoundRoute(app) // Register route for 404 Error.

// 	// Start server (with graceful shutdown).
// 	utils.StartServerWithGracefulShutdown(app)
// }

func main() {
	err := LoadEnvironmentVariables()
	if err != nil {
		log.Fatal("env loading error", err)
	}

	// 1. Configure the example database connection.
	database, err := db.CreateDB()
	if err != nil {
		log.Fatal("creating db error", err)
	}

	err = utils.CreateStocksOnDatabase(database)
	if err != nil {
		fmt.Println(err.Error())
	}

	// deneme1()

	stocks, _ := db.GetAllStocks(database)

	cookie, err := yahoo.GetCookie()
	crumb, err := yahoo.GetCrumb(cookie)

	start := time.Now()
	for i, stock := range stocks {
		if i == 1 {
			break
		}

		fmt.Println("stock:", stock)
		stock := models.Stock{Code: "EKIZ"}
		yahooChart, _ := stock.GetYahooChart("1d", "5d")
		fmt.Println("Stock:", yahooChart.Chart.Result[0].Meta.Symbol, "=", yahooChart.Chart.Result[0].Meta.RegularMarketPrice)

		yahooQuoteResponse, _ := stock.GetYahooQuoteResponse(cookie, crumb)
		fmt.Printf("QuoteResponse: %+v\n", yahooQuoteResponse.QuoteResponse.Result)

	}
	timeElapsed := time.Since(start)
	fmt.Printf("The `for` loop took %s\n", timeElapsed)

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
}
