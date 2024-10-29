package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Stock struct {
	Company string
	Price   string
	Change  string
}

func main() {
	tickers := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"WMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	stocks := []Stock{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML("section[data-testid='quote-hdr']", func(e *colly.HTMLElement) {
		stock := &Stock{}
		stock.Company = e.ChildText("h1")
		stock.Price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")
		stock.Change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']")

		stocks = append(stocks, *stock)
	})

	c.Wait()

	for _, t := range tickers {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}

  file, err := os.Create("stock.csv")
  if err != nil {
    log.Fatalln("Failed to create csv file:", err)
    os.Exit(1)
  }
  defer file.Close()

  writer := csv.NewWriter(file)
  headers := []string{
    "Company",
    "Price",
    "Change",
  }
  writer.Write(headers)
  defer writer.Flush()

  for _, stock := range stocks {
    record := []string{
      stock.Company,
      stock.Price,
      stock.Change,
    }
    writer.Write(record)
  }
}
