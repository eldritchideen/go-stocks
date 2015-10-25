package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB
var epoch time.Time

type Tick struct {
	date   string
	open   float64
	low    float64
	high   float64
	close  float64
	volume int
}

func init() {
	var _ interface{}
	epoch, _ = time.Parse("2006-01-02", "1995-01-01")
}

func check(checking string, err error) {
	if err != nil {
		log.Fatal(checking, err)
	}
}

func makeURL(stock string, startDay, startMonth, startYear, endDay, endMonth, endYear int, period string) string {
	return "http://real-chart.finance.yahoo.com/table.csv?s=" + stock +
		"&a=" + strconv.Itoa(startMonth-1) + "&b=" + strconv.Itoa(startDay) +
		"&c=" + strconv.Itoa(startYear) + "&d=" + strconv.Itoa(endMonth-1) +
		"&e=" + strconv.Itoa(endDay) + "&f=" + strconv.Itoa(endYear) +
		"&g=" + period + "&ignore=.csv"
}

func parseArgs() string {
	var fileName string

	flag.StringVar(&fileName, "file", "index.txt", "File name containing a list of stock codes to fetch")
	flag.Parse()
	return fileName
}

func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func getYahooStockData(url string) []Tick {
	ticks := make([]Tick, 0)

	res, err := http.Get(url)
	check("Error getting URL: ", err)
	webData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	check("Error reading body:", err)

	csvRows := strings.Split(string(webData), "\n")

    // Skip header row and build up tick data
	for _, r := range csvRows[1:] {
		t := strings.Split(r, ",")
		if len(t) == 7 {
			ticks = append(ticks, Tick{t[0], toFloat(t[1]), toFloat(t[3]), toFloat(t[2]),
				toFloat(t[4]), toInt(t[5])})
		}
	}
    
    // Sleep a bit so not to overload Yahoo site and get blocked. 
	time.Sleep(750 * time.Millisecond)
	return ticks
}

func updateDB(stock string, ticks []Tick) {
	stmt, err := db.Prepare("INSERT INTO ticks VALUES($1, $2, $3, $4, $5, $6, $7)")
	check("Prepare failed: ", err)

	for _, t := range ticks {
		_, err := stmt.Exec(stock, t.date, t.open, t.low, t.high, t.close, t.volume)
		check("Exec failed: ", err)
	}

}

func updateStockTicks(stock string) {

	var maxDate pq.NullTime
	var startTime time.Time

	err := db.QueryRow("select max(date)+1 from ticks where stock = $1", stock).Scan(&maxDate)
	check("Error getting max date", err)

	if maxDate.Valid {
		startTime = maxDate.Time
	} else {
		startTime = epoch
	}
	fmt.Println("startTime for", stock, "is", startTime)

	now := time.Now()

	url := makeURL(stock, int(startTime.Day()), int(startTime.Month()), startTime.Year(), int(now.Day()), int(now.Month()), now.Year(), "d")

	ticks := getYahooStockData(url)
	updateDB(stock, ticks)
}

func main() {

	fileName := parseArgs()

	fileContents, err := ioutil.ReadFile(fileName)
	check("Error reading file", err)
	stocks := strings.Split(string(fileContents), "\n")

	//Connect to DB
	db, err = sql.Open("postgres", "user=ubuntu password=ubuntu1 dbname=ubuntu sslmode=disable")
	check("Error opening DB", err)

	for _, s := range stocks {
		if s != "" {
			updateStockTicks(s)
		}
	}

	db.Close()
}
