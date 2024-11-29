package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/net/html"
)

const (
	crawlUnit      = 1
	extractionUnit = 800
	monitorUnit    = 1
	saveUnit       = 20
)

var omitedExtensions []string = []string{
	".zip", ".gz", ".msi", ".pkg",
}

func parseUrls(linkStrings []string) []*url.URL {
	links := make([]*url.URL, 0)
	for _, linkString := range linkStrings {
		if !strings.HasPrefix(linkString, "http") {
			linkString = "http://" + linkString
		}
		link, err := url.Parse(linkString)
		getBaseUrl(link)
		if err != nil {
			continue
		}
		links = append(links, link)
	}
	return links
}

func main() {
	linkToExtractC := make(chan *url.URL, 10*extractionUnit)
	linkToSaveC := make(chan *url.URL, 10*extractionUnit)
	linksC := make(chan []*url.URL, extractionUnit)
	crawlCount := make(chan int, crawlUnit)
	extractionCount := make(chan int, extractionUnit)
	crawlDone := make(chan struct{}, crawlUnit)
	done := make(chan struct{})

	baseUrls := os.Args[1:]
	parsedUrls := parseUrls(baseUrls)
	go func() {
		linksC <- parsedUrls
	}()

	for i := 0; i < monitorUnit; i++ {
		go monitor(crawlCount, extractionCount, crawlDone, linkToExtractC, linksC, done)
	}

	for i := 0; i < crawlUnit; i++ {
		fmt.Println("spawning crawl units")
		go breathFirst(linksC, linkToExtractC, linkToSaveC, crawlCount, crawlDone)
	}

	for i := 0; i < extractionUnit; i++ {
		go extract(linkToExtractC, linksC, extractionCount)
	}

	for i := 0; i < saveUnit; i++ {
		go save(linkToSaveC)
	}
	<-done
}

type seen struct {
	mu sync.Mutex
	m  map[string]struct{}
}

var seenMap = &seen{
	m: make(map[string]struct{}),
}
var baseUrl string

func breathFirst(linksC chan []*url.URL, linkToExtractC, linkToSaveC chan *url.URL, crawlCount chan int, crawlDone chan struct{}) {
	for links := range linksC {
		for _, link := range links {
			linkString := link.Host + link.Path
			if link.Path == "" {
				linkString += "/"
			}
			seenMap.mu.Lock()
			if _, ok := seenMap.m[linkString]; !ok {
				fmt.Printf("link to extract: %s\n", linkString)
				crawlCount <- 1
				seenMap.m[linkString] = struct{}{}
				linkToExtractC <- link
				linkToSaveC <- link
			}
			seenMap.mu.Unlock()
			// crawlDone <- struct{}{}
		}
	}
}

func save(linksToSave chan *url.URL) {
	for link := range linksToSave {
		resp, err := http.Get(link.String())
		if err != nil {
			continue
		}

		saveFile(link, resp)
	}
}

func saveFile(link *url.URL, resp *http.Response) {
	defer resp.Body.Close()
	dir := link.Host
	fileName := link.Path
	ext := filepath.Ext(fileName)
	if ext != "" && !unicode.IsDigit(rune(ext[len(ext)-1])) {
		index := strings.LastIndex(fileName, "/")
		var subdir = ""
		if index != -1 {
			subdir = fileName[:index]
		}
		dir = filepath.Join(dir, subdir)
		fileName = filepath.Join(dir, fileName[index+1:])
	} else {
		dir = filepath.Join(dir, fileName)
		fileName = filepath.Join(dir, "index.html")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("cannot create dir: %s", err)
	}
	f, err := os.Create(fileName)
  defer f.Close()
	if err != nil {
		log.Fatalf("cannot create file: %s", err)
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatalf("cannot write to file: %s", err)
	}
}

func getSameSiteLinks(links []*url.URL) []*url.URL {
	i := 0
	for _, link := range links {
		if link.Host == baseUrl && !omitExtenstions(link) {
			links[i] = link
			i++
		}
	}
	return links[:i]
}

func omitExtenstions(link *url.URL) bool {
	if slices.Contains(omitedExtensions, filepath.Ext(link.Path)) {
		return true
	}
	return false
}

func getBaseUrl(link *url.URL) string {
	if baseUrl == "" {
		baseUrl = link.Host
	}
	return baseUrl
}

func extractLinks(link *url.URL) []*url.URL {
	// fmt.Printf("Extracting from %s\n", link.String())
	resp, err := http.Get(link.String())
	if err != nil {
		log.Fatalf("cannot make get request: %s", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("cannot parse request body: %s", err)
	}
	links := make([]*url.URL, 0)

	f := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attrs := range n.Attr {
				if attrs.Key == "href" {
					newLink, err := resp.Request.URL.Parse(attrs.Val)
					if err != nil {
						continue
					}
					links = append(links, newLink)
				}
			}
		}
	}
	forEachNode(doc, f, nil)
	return links
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

func extract(linkToExtractC chan *url.URL, linksC chan []*url.URL, extractionCount chan int) {
	for linkToExtract := range linkToExtractC {
		links := extractLinks(linkToExtract)
		// fmt.Printf("links extracted: %s\n", linkToExtract.String())
		links = getSameSiteLinks(links)
		// for _, link := range links {
		// 	fmt.Println(link.String())
		// }
		linksC <- links
		extractionCount <- 1
	}
}

func monitor(crawlCount, extractionCount chan int, crawlDone chan struct{}, linkToExtractC chan *url.URL, linksC chan []*url.URL, done chan struct{}) {
	var count int = 0
loop:
	for {
		select {
		case <-crawlCount:
			count++
			fmt.Println("Count: ", count)
		case <-extractionCount:
			count--
			fmt.Println("Count: ", count)
		case <-crawlDone:
			if count == 0 {
				fmt.Println("Crawl done")
				close(linkToExtractC)
				close(linksC)
				break loop
			}
		}
	}
	done <- struct{}{}
}
