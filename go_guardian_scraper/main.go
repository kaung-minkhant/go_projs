package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var userAgents []string = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:130.0) Gecko/20100101 Firefox/130.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.3",
}

func selectRandomUserAgent() string {
	randoNum := rand.Int() % len(userAgents)
	return userAgents[randoNum]
}

func main() {
	rawBaseUrl := flag.String("b", "https://www.theguardian.com", "Base URL to scrape")
	flag.Parse()

	_, err := url.Parse(*rawBaseUrl)
	if err != nil {
		fmt.Printf("Cannot parse baseURL - %s: %s", *rawBaseUrl, err)
		os.Exit(1)
	}
	baseUrl := *rawBaseUrl

	// 1 the channel to hold the list of url to scrape
	worklist := make(chan []string)
	waitlist := make(chan int)

	// 2. put base url into the list
	go func() {
		worklist <- []string{baseUrl}
	}()

	// 3. seen map
	seenMap := make(map[string]bool)

	go func(worklist chan []string, waitlist chan int, seenMap map[string]bool) {
		for links := range worklist {
			for _, link := range links {
				if !seenMap[link] {
					seenMap[link] = true
          // time.Sleep(1 * time.Second)
					waitlist <- 1
					// fmt.Println("Spawning new go routine")
          fmt.Printf("Crawling URL: %s\n", link)
					go func(worklist chan []string, waitlist chan int, link string) {
						foundLinks := crawl(link, baseUrl)
						if foundLinks != nil {
							worklist <- foundLinks
							waitlist <- 0
						} else {
							waitlist <- -1
						}
					}(worklist, waitlist, link)
				}
			}
		}
	}(worklist, waitlist, seenMap)

	count := 0
	for up := range waitlist {
		if up == 1 {
			count++
		} else {
			count--
			if count == 0 && up != 0 {
				fmt.Println("Closing channels")
				close(worklist)
				close(waitlist)
			}
		}
	}

	// var i int
	// i++
	// for ; i > 0; i-- {
	// 4. listen on the channel
	// wg := &sync.WaitGroup{}
	// for links := range worklist {
	// 	// links := <-worklist
	// 	// 5. for every links
	// 	for _, link := range links {
	// 		// 6. if not seen
	// 		if !seenMap[link] {
	// 			seenMap[link] = true
	// 			// 7. crawl
	// 			// i++
	// 			wg.Add(1)
	// 			fmt.Println("Spawning one routine")
	// 			go func(link, baseUrl string, wg *sync.WaitGroup) {
	// 				fmt.Printf("Crawling URL: %s\n", link)
	// 				foundLinks := crawl(link, baseUrl)
	// 				if foundLinks != nil {
	// 					worklist <- foundLinks
	// 				} else {
	// 					wg.Done()
	// 					wg.Wait()
	// 					// fmt.Println("Closing channel")
	// 					close(worklist)
	// 				}
	// 			}(link, baseUrl, wg)
	// 		}
	// 	}
	// }
}

// func goCrawl(urls []string, baseUrl string, seenMap map[string]bool) {
// 	for _, url := range urls {
// 		if !seenMap[url] {
//       seenMap[url] = true
// 			go func(url, baseUrl string) {
// 				newUrls := crawl(url, baseUrl)
// 				goCrawl(newUrls, baseUrl, seenMap)
// 			}(url, baseUrl)
// 		}
// 	}
// }

var tokens = make(chan struct{}, 5)

func crawl(link, baseUrl string) []string {
	// return nil
	// limiting the number of requests to 5 at the same time with buffered channel
	tokens <- struct{}{}

	resp, _ := getRequest(link)

	<-tokens
	links := discoverLinks(resp, baseUrl)
	foundLinks := []string{}

	for _, link := range links {
		fullUrl := resolveRelativeLink(link, baseUrl)
		if checkHost(fullUrl, baseUrl) {
			foundLinks = append(foundLinks, fullUrl)
		}
	}

	parseHTML(resp)

	return foundLinks
}

func parseHTML(resp *http.Response) string {
	return ""
}

func getRequest(link string) (*http.Response, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", selectRandomUserAgent())

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func resolveRelativeLink(link, baseUrl string) string {
	if strings.HasPrefix(link, "/") {
		return fmt.Sprintf("%s%s", baseUrl, link)
	}
	return link
}

func discoverLinks(resp *http.Response, baseUrl string) []string {
	if resp == nil {
		return []string{}
	}

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Cannot create document from response: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if document == nil {
		return []string{}
	}

	foundUrls := []string{}

	document.Find("a").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("href")
		foundUrls = append(foundUrls, value)
	})

	return foundUrls
}

func checkHost(link, baseUrl string) bool {
	url, err := url.Parse(link)
	if err != nil {
		fmt.Printf("Cannot parse URL - %s: %s\n", link, err.Error())
		return false
	}
	base, _ := url.Parse(baseUrl)
	if url.Host != base.Host {
		return false
	}
	return true
}
