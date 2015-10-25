package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func main() {
	db, err := sql.Open("postgres", "user=ubuntu password=ubuntu1 dbname=ubuntu sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Ping failed: ", err)
	}

	stmt, err := db.Prepare("INSERT INTO ticks VALUES($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		log.Fatal("Prepare failed: ", err)
	}
	
	t, err := time.Parse("2006-01-02", "2015-09-25")
	if err != nil {
		log.Fatal("Parse date failed: ",err)
	}
	fmt.Println("Parsed date is:", t)
	// can also Exec the statement with a date column in string format of "2015-09-28" gets converted under the covers. 
	res, err := stmt.Exec("FMG.AX", t , 1.78, 1.74, 1.80, 1.80, 25927000)
	if err != nil {
		log.Fatal("Exec failed: ",err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rows affected = %d\n", rowCnt)
	db.Close()
}
