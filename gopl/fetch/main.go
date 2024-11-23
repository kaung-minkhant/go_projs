package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func checkHttp(url string) string {
  if !strings.HasPrefix(url, "http://") {
    return "http://"+url
  } 
  return url
}

func main() {
	for _, _url := range os.Args[1:] {
    url := checkHttp(_url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetching %s failed %v\n", url, err)
			return
		}
    // bytes, err := io.ReadAll(resp.Body)
    _, err = io.Copy(os.Stdout, resp.Body) // better than above because, readAll tries to allocate the buffer of response body size, copy, directly works on stream, no buffer needed
    if err != nil {
      fmt.Fprintf(os.Stderr, "Parsing body for %s falied: %v with status code %s\n", url, err, resp.Status)
      return
    }
    resp.Body.Close()
    // fmt.Fprintf(os.Stdout, " with status code %s\n", resp.Status)

    // fmt.Printf("Response for %s: %s\n", url, bytes)
	}
}
