package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type GenericResult struct {
	Page         int32  `json:"page"`
	PerPage      int32  `json:"per_page"`
	TotalResults int32  `json:"total_results"`
	PreviousPage string `json:"prev_page"`
	NextPage     string `json:"next_page"`
}

type SearchResult struct {
	*GenericResult
	Photos []*Photo `json:"photos"`
}

type CuratedResult struct {
	*GenericResult
	Photos []*Photo `json:"photos"`
}

type SearchVideoResult struct {
	*GenericResult
	Videos []*Video `json:"videos"`
}

type PopularVideosResult struct {
  *GenericResult
  Url string `json:"url"`
	Videos []*Video `json:"videos"`
}


type Photo struct {
	Id     int32  `json:"id"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Url    string `json:"url"`
}

type Video struct {
	Id       int32  `json:"id"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Url      string `json:"url"`
	Duration int32  `json:"duration"`
}

type apiConfig struct {
	ApiKey string `json:"apiKey"`
}

type apiClient struct {
	*apiConfig
	c              http.Client
	RemainingTimes int32
}

func (c *apiClient) requestWithAuth(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.ApiKey)
	resp, err := c.c.Do(req)

	remainingRequests, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return resp, nil
	}

	c.RemainingTimes = int32(remainingRequests)

	return resp, err
}

func (c *apiClient) SearchPhotos(query string, perPage int32, page int32) (*SearchResult, error) {
	search := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)

	resp, err := c.requestWithAuth("GET", search)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResult SearchResult
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, err
	}
	return &searchResult, nil
}

func (c *apiClient) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?page=%d&per_page=%d", page, perPage)

	resp, err := c.requestWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var curatedResult CuratedResult
	if err := json.Unmarshal(body, &curatedResult); err != nil {
		return nil, err
	}
	return &curatedResult, nil
}

func (c *apiClient) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)
	resp, err := c.requestWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var photo Photo
	if err := json.Unmarshal(body, &photo); err != nil {
		return nil, err
	}

	return &photo, nil
}

func getRandomId() int {
	randomNumber := rand.Intn(1001)
	return randomNumber
}

func (c *apiClient) GetRandomPhoto() (*Photo, error) {
	id := getRandomId()
	result, err := c.CuratedPhotos(1, id)
	if err != nil {
		return nil, err
	}

	if len(result.Photos) == 1 {
		return result.Photos[0], nil
	}
	return &Photo{}, nil
}

func (c *apiClient) SearchVideo(query string, page, per_page int32) (*SearchVideoResult, error) {
  url := fmt.Sprintf(VideoApi+"/search?query=%s&page=%d&per_page=%d", query, page, per_page)  
  resp, err := c.requestWithAuth("GET", url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body);
  if err != nil {
    return nil, err
  }
  var searchResult SearchVideoResult
  if err := json.Unmarshal(body, &searchResult); err != nil {
    return nil, err
  }
  return &searchResult, nil
}


func (c *apiClient) GetPopularVideo( page, per_page int32) (*PopularVideosResult, error) {
  url := fmt.Sprintf(VideoApi+"/popular?page=%d&per_page=%d", page, per_page)  
  resp, err := c.requestWithAuth("GET", url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body);
  if err != nil {
    return nil, err
  }
  var searchResult PopularVideosResult
  if err := json.Unmarshal(body, &searchResult); err != nil {
    return nil, err
  }
  return &searchResult, nil
}

func (c *apiClient) GetVideo(id int32) (*Video, error) {
  url := fmt.Sprintf(VideoApi+"/video/%d", id)  
  resp, err := c.requestWithAuth("GET", url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body);
  if err != nil {
    return nil, err
  }
  var video Video 
  if err := json.Unmarshal(body, &video); err != nil {
    return nil, err
  }
  return &video, nil
}

func (c *apiClient) GetRandomVideo() (*Video, error) {
  id := getRandomId()
  videosResult, err := c.GetPopularVideo(int32(id), 1)
  if err != nil {
    return nil, err
  }
  if len(videosResult.Videos) == 1 {
    return videosResult.Videos[0], nil
  }
  return &Video{}, nil
}

func NewClient(token string) *apiClient {
	c := http.Client{}

	return &apiClient{
		apiConfig: &apiConfig{
			ApiKey: token,
		},
		c: c,
	}
}

var api apiConfig

func loadConfig(filename string) {
	body, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Cannot read config file", err)
		panic("Error")
	}
	if err := json.Unmarshal(body, &api); err != nil {
		log.Fatal("Cannot unmarshal config file", err)
		panic("Error")
	}

}

func main() {
	loadConfig(".config")
	c := NewClient(api.ApiKey)
	// photo, err := c.GetRandomPhoto()
	// photos, err := c.CuratedPhotos(5, 1)
  video, err := c.GetRandomVideo()
	if err != nil {
		log.Fatal("Getting photo error", err)
		return
	}

	result, _ := json.Marshal(video)
	fmt.Printf("Photo(s): %s\n", result)
	fmt.Printf("Remining Requests: %d", c.RemainingTimes)
}
