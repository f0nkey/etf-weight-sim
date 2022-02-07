package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Security struct {
	Name string `json:"name"`
	Ticker string `json:"ticker"`
	Weight float64 `json:"weight"`
	DailyChange float64 `json:"daily_change"`
}

func main() {
	link := `https://www.ssga.com/us/en/individual/etfs/library-content/products/fund-data/etfs/us/holdings-daily-us-en-spy.xlsx`
	securities := spdrETFSecurityHoldings(link)

	var tickers []string
	for _, security := range securities {
		tickers = append(tickers, security.Ticker)
	}
	marketChanges, err := tickerMarketChanges(tickers)
	if err != nil {
		log.Fatalln(err)
	}
	for i, security := range securities {
		ch := marketChanges[security.Ticker]
		security.DailyChange = ch
		securities[i] = security
	}

	b, err := json.MarshalIndent(&securities,"","\t")
	if err != nil {
		log.Fatalln(err)
	}

	if err = ioutil.WriteFile("../weight-sim/src/assets/spy.json",b,0644); err != nil {
		log.Fatalln(err)
	}
}

func spdrETFSecurityHoldings(url string) []Security {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	f, err := excelize.OpenReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var securities []Security

	rows, err := f.GetRows("holdings")
	if err != nil {
		log.Fatalln(err)
	}
	for _, row := range rows {
		if len(row) == 8 && row[1] != "Ticker" {
			w, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				log.Fatalln(err)
			}
			securities = append(securities, Security{
				Name:   row[0],
				Ticker: row[1],
				Weight: w,
			})
		}
	}
	return securities
}

func tickerMarketChanges(tickers []string) (map[string]float64, error) {
	url := fmt.Sprintf(`https://query1.finance.yahoo.com/v7/finance/quote?symbols=%sfields=regularMarketChangePercent`, strings.Join(tickers, ","))
	cl := http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := cl.Get(url)
	if err != nil {
		return nil, err
	}

	marketData := struct {
		QuoteResponse struct {
			Result []struct {
				RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
				Symbol                     string  `json:"symbol"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"quoteResponse"`
	}{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &marketData)
	if err != nil {
		return nil, err
	}

	marketChanges := make(map[string]float64)
	for _, quote := range marketData.QuoteResponse.Result {
		marketChanges[quote.Symbol] = quote.RegularMarketChangePercent
	}
	return marketChanges, nil
}
