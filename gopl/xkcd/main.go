package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Comic struct {
	Num        int
	Day        string
	Month      string
	Year       string
	SafeTitle  string `json:"safe_title"`
	Title      string
	Transcript string
	Img        string
	Alt        string
}

const (
	API_URL          = "https://xkcd.com/"
	NumberOfFetchers = 200
	NumberOfIndexers = 20
)

func doGetRequest(url string) (*http.Response, error) {
	client := http.Client{}
	newGetRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create new get request: %s", err)
	}
	resp, err := client.Do(newGetRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot make get request: %s", err)
	}
	return resp, nil
}

func getNumberOfComic() (int, error) {
	url := fmt.Sprintf("%s/info.0.json", API_URL)
	resp, err := doGetRequest(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return 0, fmt.Errorf("cannot decode body for comic number: %s", err)
	}
	return comic.Num, nil
}

func get(n int) (*Comic, error) {
	url := fmt.Sprintf("%s/%d/info.0.json", API_URL, n)
	resp, err := doGetRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var comic = new(Comic)
	if err := json.NewDecoder(resp.Body).Decode(comic); err != nil {
		return nil, fmt.Errorf("cannot decode body for comic: %s", err)
	}
	return comic, nil
}

// get number of comics
// spawn works and provide the channels
// spawn indexing workers

type indexData struct {
	mu        sync.Mutex
	wordIndex map[string]map[int]bool
	numIndex  map[int]Comic
}

// TODO: change this into io.Writer
func index(fileName string) {
	maxNumberOfComic, err := getNumberOfComic()
	if err != nil {
		log.Fatal(err)
	}
	// maxNumberOfComic = 50

	// comics := make(chan *Comic, maxNumberOfComic)
	comics := make(chan *Comic, NumberOfFetchers*2)
	fetcherDone := make(chan struct{})
	indexDone := make(chan struct{})

	data := &indexData{}
	data.wordIndex = make(map[string]map[int]bool)
	data.numIndex = make(map[int]Comic)

	spawnComicFetchers(maxNumberOfComic, comics, fetcherDone)
	spawnIndexWorkers(comics, indexDone, data)

	go func() {
		for i := 0; i < NumberOfFetchers; i++ {
			<-fetcherDone
		}
		close(comics)
		close(fetcherDone)
	}()
	for i := 0; i < NumberOfIndexers; i++ {
		<-indexDone
	}
	close(indexDone)
	encodeIntoFile(fileName, data)
}

func encodeIntoFile(fileName string, data *indexData) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("cannot create index file: %s", err)
	}
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(data.numIndex); err != nil {
		return fmt.Errorf("cannot encode numIndex: %s", err)
	}
	if err := encoder.Encode(data.wordIndex); err != nil {
		return fmt.Errorf("cannot encode wordIndex: %s", err)
	}
	return nil
}

func decodeFromFile(fileName string) (*indexData, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open index file: %s", err)
	}
	decoder := gob.NewDecoder(f)
	data := new(indexData)
	if err := decoder.Decode(&data.numIndex); err != nil {
		return nil, fmt.Errorf("cannot decode numIndex: %s", err)
	}
	if err := decoder.Decode(&data.wordIndex); err != nil {
		return nil, fmt.Errorf("cannot decode wordIndex: %s", err)
	}
	return data, nil
}

func spawnIndexWorkers(comics chan *Comic, done chan struct{}, data *indexData) {
	fmt.Println("spawning index workers")
	for i := 0; i < NumberOfIndexers; i++ {
		go indexWorker(comics, done, data)
	}
}

func indexWorker(comics chan *Comic, done chan struct{}, data *indexData) {
	for comic := range comics {
		scanner := bufio.NewScanner(strings.NewReader(comic.Transcript))
		scanner.Split(bufio.ScanWords)
		data.mu.Lock()
		for scanner.Scan() {
			word := strings.ToLower(scanner.Text())
			// wordS := make([]rune, len(word))
			// for _, char := range word {
			// 	if !unicode.IsPunct(char) {
			// 		wordS = append(wordS, char)
			// 	}
			// }
			// word = string(wordS)
			// fmt.Println(word)
			data.numIndex[comic.Num] = *comic
			if _, ok := data.wordIndex[word]; !ok {
				data.wordIndex[word] = make(map[int]bool)
			}
			data.wordIndex[word][comic.Num] = true
		}
		data.mu.Unlock()
	}
	done <- struct{}{}
}

func spawnComicFetchers(maxNumberOfComic int, comics chan *Comic, done chan struct{}) {
	comicNumbers := make(chan int, maxNumberOfComic)

	fmt.Println("spawning fetch workers")
	for i := 0; i < NumberOfFetchers; i++ {
		go fetcher(comicNumbers, comics, done)
	}
	for i := 1; i <= maxNumberOfComic; i++ {
		if i == 404 {
			continue
		}
		comicNumbers <- i
	}
	close(comicNumbers)
}

func fetcher(comicNumber chan int, comics chan *Comic, done chan struct{}) {
	for n := range comicNumber {
		url := fmt.Sprintf("%s/%d/info.0.json", API_URL, n)
		resp, err := doGetRequest(url)
		if err != nil {
			fmt.Printf("Error fetching comic number %d: %s\n", n, err)
			continue
		}
		var comic = new(Comic)
		if err := json.NewDecoder(resp.Body).Decode(comic); err != nil {
			fmt.Printf("Error decoding comic number %d: %s\n", n, err)
			resp.Body.Close()
			continue
		}
		comics <- comic
	}
	done <- struct{}{}
}

func getComicContainingQuery(query []string, data *indexData) []Comic {
	// make a count to keep track of match hit
	hit := make(map[int]int)
	selectedComics := []Comic{}
	for _, word := range query {
		lower := strings.ToLower(word)
		if _, ok := data.wordIndex[lower]; ok {
			for comicNumber := range data.wordIndex[lower] {
				hit[comicNumber]++
			}
		}
	}
	for comicNum, hitCount := range hit {
		if hitCount == len(query) {
			selectedComics = append(selectedComics, data.numIndex[comicNum])
		}
	}
	return selectedComics
}

func search(query []string, fileName string) []Comic {
	data, err := decodeFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("data: %#v\n", data)
	comics := getComicContainingQuery(query, data)
	return comics
}

func main() {
	// index("index.index")
  comics := search([]string{"thanks!"}, "index.index")
	fmt.Printf("Comics: %#v\n", comics)
}
