package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func max(highs []float64) float64 {
	max := highs[0]
	for _, h := range highs[1:] {
		if h > max {
			max = h
		}
	}
	return max
}

func main() {
	db, err := sql.Open("postgres", "user=ubuntu password=ubuntu1 dbname=ubuntu sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("select stock, high from ticks where date >= (current_date - 365) order by date desc")
	if err != nil {
		log.Fatal("Query failed: ", err)
	}

	var stock string
	var high float64

	highs := make(map[string][]float64)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&stock, &high)
		if err != nil {
			log.Fatal(err)
		}
		highs[stock] = append(highs[stock], high)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Searching for stocks that have made a new 12 month high.")
	
	for stock,highs := range highs {
		currentHigh := highs[0]
		if currentHigh > max(highs[1:]) {
			log.Println(stock, "made a new high of", currentHigh)
		}
	} 
	log.Println("Done.")
}
