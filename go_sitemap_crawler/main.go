package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type SEOData struct {
	URL             string
	Title           string
	H1              string
	MetaDescription string
	StatusCode      int
}

type parser interface {
	getSEOData(*http.Response) (SEOData, error)
}

type DefaultParser struct{}

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

func scrapeSiteMap(url string, p parser, concurrency int) []SEOData {
	pagesToCrawl := exploreSitemapsAndGetPages(url)
	res := scrapePages(pagesToCrawl, p, concurrency)
	return res
}

func (d *DefaultParser) getSEOData(resp *http.Response) (SEOData, error) {
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Cannot create document: %s\n", err)
		return SEOData{}, nil
	}
	resp.Body.Close()
	result := &SEOData{}
	result.URL = resp.Request.URL.String()
	result.StatusCode = resp.StatusCode
	result.H1 = document.Find("h1").First().Text()
	result.MetaDescription, _ = document.Find("meta[^=description]").Attr("content")
	result.Title = document.Find("title").First().Text()
	return *result, nil
}

func scrapePages(urls []string, p parser, concurrency int) []SEOData {
	seenPages := make(map[string]bool)
	results := []SEOData{}
	worklist := make(chan string)
	resultList := make(chan SEOData)
	wg := &sync.WaitGroup{}

	go func(urls []string, seenPages map[string]bool, worklist chan string, resultList chan SEOData, wg *sync.WaitGroup) {
		for _, url := range urls {
			if !seenPages[url] {
				seenPages[url] = true
				worklist <- url
			}
		}
		wg.Wait()
		close(worklist)
		close(resultList)
	}(urls, seenPages, worklist, resultList, wg)

  tokens := make(chan struct{}, concurrency)
	go func(worklist chan string, resultList chan SEOData, p parser, wg *sync.WaitGroup, tokens chan struct{}) {
		for link := range worklist {
			wg.Add(1)
      tokens <- struct{}{}
			go func(link string, p parser, resultList chan SEOData, wg *sync.WaitGroup, tokens chan struct{}) {
				defer wg.Done()
        fmt.Printf("Scrapping %s\n", link)
				result, err := scrapePage(link, p)
        <- tokens
				if err != nil {
					fmt.Printf("Error scraping page: %s: %s\n", link, err)
				} else {
					resultList <- *result
				}
			}(link, p, resultList, wg, tokens)
		}
	}(worklist, resultList, p, wg, tokens)

	for result := range resultList {
		results = append(results, result)
	}

	return results
}

func scrapePage(url string, p parser) (*SEOData, error) {
	resp, err := makeRequest(url)
	if err != nil {
		fmt.Printf("Request failed to %s: %s\n", url, err)
		return nil, err
	}
	result, err := p.getSEOData(resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// it takes a link,
// go to that link, get all urls,
// seperate them between sitemaps and pages
// add pages to crawl later
// go to each sitemap, and get more sitemaps (if available) and pages
// and repeat
// return pages to crawl

func exploreSitemapsAndGetPages(url string) []string {
	fmt.Printf("Exploring sitemap: %s\n", url)
	seenSitemaps := make(map[string]bool)
	worklist := make(chan []string)
  pagesChan := make(chan []string)
  statusC := make(chan int)

	pagesToCrawl := []string{}

	go func() { worklist <- []string{url} }()

	go func(worklist chan []string, pagesC chan []string, statusC chan int) {
		for links := range worklist {
			for _, link := range links {
				if !seenSitemaps[link] {
					seenSitemaps[link] = true
          statusC <- 1
					go func(link string, worklist chan []string, pagesC chan []string, statusC chan int) {
            fmt.Printf("Requesting %s\n", link)
						resp, _ := makeRequest(link)
						urls := extractUrls(resp)
						sitemaps, pages := seperateSitmapsFromPages(urls)
						if len(sitemaps) != 0 {
							worklist <- sitemaps
              statusC <- 0
						} else {
              statusC <- -1
            }
            pagesChan <- pages
					}(link, worklist, pagesChan, statusC)
				}
			}
		}
	}(worklist, pagesChan, statusC)

  go func (statusC chan int)  {
    counter := 0
    for status := range statusC {
      if status == 1 {
        counter++
      } else {
        counter--
      }
      if counter == 0 && status == -1 {
        close(statusC)
        close(worklist)
        close(pagesChan)
      }
    } 
  }(statusC)

  for pages := range pagesChan {
    pagesToCrawl = append(pagesToCrawl, pages...)
  }

	return pagesToCrawl
}

func seperateSitmapsFromPages(urls []string) ([]string, []string) {
	sitemaps := []string{}
	pages := []string{}
	for _, url := range urls {
		if strings.Contains(url, ".xml") {
			sitemaps = append(sitemaps, url)
		} else {
			pages = append(pages, url)
		}
	}
	return sitemaps, pages
}

func extractUrls(resp *http.Response) []string {
	document, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Cannot create document from response: %s\n", err)
		return []string{}
	}

	var urls []string
	document.Find("loc").Each(func(i int, s *goquery.Selection) {
		url := s.Text()
		urls = append(urls, url)
	})
	return urls
}

func makeRequest(url string) (*http.Response, error) {
	client := http.Client{}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", getRandomUserAgent())

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	p := &DefaultParser{}
	results := scrapeSiteMap("https://www.quicksprout.com/sitemap.xml", p, 10)
	for _, result := range results {
		fmt.Println(result)
	}
}
