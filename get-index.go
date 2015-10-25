package main

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)	
	}
}

func getHTML() {
	var allOrds string
	var count = 0
	
	res, err := http.Get("http://www.marketindex.com.au/all-ordinaries")
	check(err)
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	check(err)
	
	log.Println("Extracting shares from html.")
	doc.Find("#asx_sp_table > tbody > tr > td:nth-child(2)").Each(
		func(i int, s *goquery.Selection) {
			count += 1
			allOrds += s.Text() + ".AX\n"
		})
	log.Println("Writing", count, "shares to file index.txt")	
	ioutil.WriteFile("index.txt", []byte(allOrds), 0644)
}

func main() {
	getHTML()
}
