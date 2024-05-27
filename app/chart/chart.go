package chart

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cinar/indicator"
	"github.com/harunalbayrak/go-finance-app/app/models"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func StockChart(stock string, last int, interval int) {
	var xvalues []time.Time
	var yvalues []float64
	var volumevalues []float64

	stockData, _ := readStocksData(stock)

	for i := 0; i < len(stockData.Prices); i = i + interval {
		parsed, _ := time.Parse(chart.DefaultDateFormat, stockData.Prices[i].Date)
		xvalues = append(xvalues, parsed)
		yvalues = append(yvalues, stockData.Prices[i].Value)
		volumevalues = append(volumevalues, stockData.Prices[i].Volume/5000)
	}

	var xvalues2 []time.Time
	var yvalues2 []float64
	var volumevalues2 []float64

	if len(xvalues) > last {
		xvalues2 = xvalues[len(xvalues)-last:]
		yvalues2 = yvalues[len(yvalues)-last:]
		volumevalues2 = volumevalues[len(volumevalues)-last:]
	} else {
		xvalues2 = xvalues
		yvalues2 = yvalues
		volumevalues2 = volumevalues
	}

	result := indicator.Sma(20, yvalues)
	ema12 := indicator.Ema(12, yvalues)
	ema26 := indicator.Ema(26, yvalues)
	macd, signal := indicator.Macd(yvalues)
	rs, rsi := indicator.Rsi(yvalues)
	rs2, rsi2 := indicator.Rsi2(yvalues)
	fmt.Printf("SMA: %.2f\n", result[len(result)-1])
	fmt.Printf("EMA12: %.2f\n", ema12[len(ema12)-1])
	fmt.Printf("EMA26: %.2f\n", ema26[len(ema26)-1])
	fmt.Printf("Macd: %.2f\n", macd[len(macd)-1])
	fmt.Printf("Signal: %.2f\n", signal[len(signal)-1])
	fmt.Printf("RS: %.2f\n", rs[len(rs)-1])
	fmt.Printf("RSI: %.2f\n", rsi[len(rsi)-1])
	fmt.Printf("RS2: %.2f\n", rs2[len(rs2)-1])
	fmt.Printf("RSI2: %.2f\n", rsi2[len(rsi2)-1])

	priceSeries := chart.TimeSeries{
		Name: "EKIZ",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorBlue,
		},
		XValues: xvalues2,
		YValues: yvalues2,
	}

	smaSeries := chart.SMASeries{
		Name: "SMA 16",
		Style: chart.Style{
			Show:            true,
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		Period:      16,
		InnerSeries: priceSeries,
	}

	bbSeries := &chart.BollingerBandsSeries{
		Name: "Bol. Bands",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("efefef"),
			FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
		},
		InnerSeries: priceSeries,
	}

	emaSeries := &chart.EMASeries{
		Name: "EMA 12",
		Style: chart.Style{
			Show:            true,
			StrokeColor:     drawing.ColorFromHex("782482"),
			StrokeDashArray: []float64{5.0, 5.0},
		},
		Period:      12,
		InnerSeries: priceSeries,
	}

	volumeSeries := chart.TimeSeries{
		Name: "Volume",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorGreen,
			FillColor:   chart.ColorGreen.WithAlpha(64),
		},
		XValues: xvalues2,
		YValues: volumevalues2,
	}

	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		XAxis: chart.XAxis{
			Name:           "Date",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(),
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			Name:      "Price",
			// Range: &chart.ContinuousRange{
			// 	Max: 200.0,
			// 	Min: 0.0,
			// },
		},
		Series: []chart.Series{
			priceSeries,
			smaSeries,
			emaSeries,
			bbSeries,
			volumeSeries,
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, _ := os.Create("charts/" + stock + ".png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func readStocksData(stock string) (*models.StockC, error) {
	file, err := os.Open("data/" + stock + "-1d.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var stockData models.StockC
	var prices []models.StockPrice
	for i, record := range records {
		fmt.Println(record)

		if i == 0 {
			continue // başlık satırını atla
		}

		value, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		volume, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		price := models.StockPrice{
			Date:   record[0],
			Value:  value,
			Volume: volume,
		}
		prices = append(prices, price)
	}

	stockData = models.StockC{
		Code:   stock,
		Prices: prices,
	}

	return &stockData, nil
}
