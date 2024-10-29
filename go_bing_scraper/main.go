package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:130.0) Gecko/20100101 Firefox/130.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.3",
}

func getRandomUserAgent() string {
	randomNum := rand.Int() % len(userAgents)
	return userAgents[randomNum]
}

var bingDomains = map[string]string{
	"com": "",
}

type SearchResult struct {
	// ResultRank  int
	ResultUrl   string
	ResultTitle string
	ResultDesc  string
}

// https://bing.com/search?q="dog"&count=10
// https://www.bing.com/search?q=dog
// https://www.bing.com/search?q=dog&first=11&count=10&cc=au
func buildBingUrls(query string, page, count int, country string) []string {
	results := []string{}
	cc, exist := bingDomains[country]
	if !exist {
		fmt.Printf("Country %s is not support!\n", country)
		return results
	}

	for i := 0; i < page; i++ {
		first := getStartingItemNumber(i, count)
		result := fmt.Sprintf("https://www.bing.com/search?q=%s&first=%d&count=%d&cc=%s", query, first, count, cc)
		results = append(results, result)
	}
	return results
}

func getStartingItemNumber(page, count int) int {
	return page*count + 1
}

func makeRequest(link string, proxy interface{}) (*http.Response, error) {
	// in case there is a proxy server to use
	var client *http.Client
	switch V := proxy.(type) {
	case string:
		proxyUrl, _ := url.Parse(V)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: func(r *http.Request) (*url.URL, error) {
					return proxyUrl, nil
				},
			},
		}
	default:
		client = &http.Client{}
	}

	newRequest, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	newRequest.Header.Add("User-Agent", getRandomUserAgent())
	res, err := client.Do(newRequest)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func scrapePage(page string, proxy interface{}) ([]SearchResult, error) {
	result := []SearchResult{}
	resp, err := makeRequest(page, proxy)
	if err != nil {
		fmt.Printf("Cannot scrape page %s: %s\n", page, err)
		return nil, err
	}
	result = parseBingResult(resp)
	return result, nil
}

const concurrency = 3

func scrapePages(pages []string, proxy interface{}) ([]SearchResult, error) {
	tokensC := make(chan struct{}, concurrency)
	resultC := make(chan []SearchResult)
	wg := &sync.WaitGroup{}
	results := []SearchResult{}

	go func(pages []string, wg *sync.WaitGroup, tokenC chan struct{}, resultC chan []SearchResult, proxy interface{}) {
		for _, page := range pages {
			wg.Add(1)
			tokensC <- struct{}{}
			go func(page string, wg *sync.WaitGroup, tokensC chan struct{}, resultC chan []SearchResult, proxy interface{}) {
				fmt.Printf("Scrapping page %s\n", page)
				defer wg.Done()
				result, err := scrapePage(page, proxy)
				if err != nil || result == nil {
					<-tokensC
					return
				} else {
					<-tokensC
					resultC <- result
				}

			}(page, wg, tokensC, resultC, proxy)
		}
		wg.Wait()
		close(resultC)
		close(tokenC)
	}(pages, wg, tokensC, resultC, proxy)

	for result := range resultC {
		results = append(results, result...)
	}
	return results, nil
}

func parseBingResult(resp *http.Response) []SearchResult {
	results := []SearchResult{}
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Cannot create new document: %s\n", err)
		return nil
	}

	sel := document.Find("li.b_algo")
	for i := range sel.Nodes {
		var result SearchResult
		item := sel.Eq(i)
		link, ok := item.Find("a").Attr("href")
		if !ok {
			continue
		}
		if strings.HasPrefix(link, "/") || link == "#" {
			continue
		}
		result.ResultUrl = link
		title := item.Find("h2").Text()
		result.ResultTitle = title
		description := item.Find("p").Text()
		result.ResultDesc = description
		results = append(results, result)
	}
	return results
}

func ScrapeBing(query, domain string, page, count int, proxy interface{}) ([]SearchResult, error) {
	bingUrls := buildBingUrls(query, page, count, domain)

	results, err := scrapePages(bingUrls, proxy)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}
	return results, nil
}

func main() {
  query := flag.String("q", "dog", "Query to search on bing")
  flag.Parse()

	res, err := ScrapeBing(*query, "com", 100, 50, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Result count:", len(res))
	for _, result := range res {
		fmt.Printf("Result: %+v\n", result.ResultUrl)
	}
}
