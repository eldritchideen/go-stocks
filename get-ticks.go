package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type stockDetails struct {
	stockCode  string
	startDay   int
	startMonth int
	startYear  int
	endDay     int
	endMonth   int
	endYear    int
	period     string
}

func makeURL(stock *stockDetails) string {
	return "http://real-chart.finance.yahoo.com/table.csv?s=" + stock.stockCode +
		"&a=" + strconv.Itoa(stock.startMonth-1) + "&b=" + strconv.Itoa(stock.startDay) +
		"&c=" + strconv.Itoa(stock.startYear) + "&d=" + strconv.Itoa(stock.endMonth-1) +
		"&e=" + strconv.Itoa(stock.endDay) + "&f=" + strconv.Itoa(stock.endYear) +
		"&g=" + stock.period + "&ignore=.csv"
}

func parseArgs() stockDetails {

	stock := stockDetails{}

	flag.StringVar(&stock.stockCode, "code", "", "Yahoo finance stock code")
	flag.IntVar(&stock.startDay, "sd", 1, "Start Day")
	flag.IntVar(&stock.startMonth, "sm", 1, "Start Month")
	flag.IntVar(&stock.startYear, "sy", 1980, "Start Year")
	flag.IntVar(&stock.endDay, "ed", 0, "End Day")
	flag.IntVar(&stock.endMonth, "em", 0, "End Month")
	flag.IntVar(&stock.endYear, "ey", 0, "End Year")
	flag.StringVar(&stock.period, "p", "d", "period of prices")

	flag.Parse()

	if stock.endDay == 0 || stock.endMonth == 0 || stock.endYear == 0 {
		t := time.Now()
		stock.endDay = int(t.Day())
		stock.endMonth = int(t.Month())
		stock.endYear = t.Year()
	}

	return stock
}

func getStockData(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Error getting URL: ", err)
	}
	csv, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(csv)
}

func main() {
	s := parseArgs()
	url := makeURL(&s)
	fmt.Println(getStockData(url))
}
