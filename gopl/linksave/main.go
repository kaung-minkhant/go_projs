package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

var baseUrl string = ""

func getSameSite(links []string) []string {
	i := 0
	if baseUrl != "" {
		for _, link := range links {
			_url, err := url.Parse(link)
			if err != nil {
				continue
			}
			if _url.Host == baseUrl {
				links[i] = link
				i++
			}
		}
	}
	return links[:i]
}

func extract(link string) ([]string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)

	getLink := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url, err := resp.Request.URL.Parse(attr.Val)
					if err != nil {
						continue
					}
					if url.Host == baseUrl {
						link := url.String()
						if strings.HasSuffix(link, "/") {
							link = link[:len(link)-1]
						}
						links = append(links, link)
					}
				}
			}
		}
	}

	forEachNode(doc, getLink, nil)

	return links, nil
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

func parseRequestUrl(link string) (*url.URL, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func save(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get request failed: %s", err)
	}
	full, err := parseRequestUrl(url)
	if baseUrl == "" {
		if err != nil {
			return fmt.Errorf("parsing url failed: %s", err)
		}
		baseUrl = full.Host
	}
	defer resp.Body.Close()
	dir := full.Host
	fileName := full.Path
	ext := filepath.Ext(fileName)
	if ext == "" {
		dir = filepath.Join(dir, fileName)
	  fileName = filepath.Join(dir, "index.html")
	} else {
		slash := strings.LastIndex(fileName, "/")
		dir = filepath.Join(dir, fileName[:slash])
    fileName = filepath.Join(dir, fileName[slash:])
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("cannot create directory: %s", err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("cannot create file: %s", err)
	}
  defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("cannot write to file: %s", err)
	}

	return nil
}

func checkHttp(url string) string {
	if !strings.HasPrefix(url, "http") {
		return "http://" + url
	}
	return url
}

func crawl(link string) []string {
	link = checkHttp(link)
	fmt.Println(link)
	err := save(link)
	if err != nil {
		log.Fatalf("Cannot save file: %s\n", err)
	}
	links, err := extract(link)
	links = getSameSite(links)
	if err != nil {
		log.Println(err)
		return nil
	}
	return links
}

func breathFirst(links []string, f func(url string) []string, pre func(url string) string) {
	seen := make(map[string]struct{})
	for len(links) > 0 {
		workList := links
		links = nil
		for _, item := range workList {
			item = pre(item)
			if _, ok := seen[item]; !ok {
				seen[item] = struct{}{}
				links = append(links, f(item)...)
			}
		}
	}
}


// 
func main() {
	f := func(link string) string {
		parsed, err := url.Parse(link)
		if err != nil {
			log.Fatal(err)
		}
		if parsed.Path == "/" {
			return parsed.Host
		}
		return parsed.Host + parsed.Path
	}
	breathFirst(os.Args[1:], crawl, f)
}
