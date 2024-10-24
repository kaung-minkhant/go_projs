package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type geoLocationData struct {
	Name      string  `json:"name"`
	Lattitude float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
}

type geoLocationArr []geoLocationData

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (*apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config apiConfigData

	if err = json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}

func getGeolocation(city string) (string, string, error) {
	config, err := loadApiConfig("./.apiConfig")
	if err != nil {
		log.Fatal("Could not read config file", err)
		return "", "", err
	}
	resp, err := http.Get("http://api.openweathermap.org/geo/1.0/direct?q=" + city + "&appid=" + config.OpenWeatherMapApiKey)
	if err != nil {
		log.Fatal("Could not get geolocation api", err)
		return "", "", err
	}
	defer resp.Body.Close()
  var geoLocations geoLocationArr
  if err := json.NewDecoder(resp.Body).Decode(&geoLocations); err != nil {
		log.Fatal("Could not parse body of geolocation api", err)
    return "", "", err
  }
  lat := strconv.FormatFloat(geoLocations[0].Lattitude, 'g', -1, 64)
  long := strconv.FormatFloat(geoLocations[0].Longitude, 'g', -1, 64)

	return lat, long, nil
}

func getWeather(city string) (*weatherData, error) {
	config, err := loadApiConfig("./.apiConfig")
	if err != nil {
		log.Fatal("Could not read config file", err)
		return nil, err
	}

  lat, long, err := getGeolocation(city)
  if err != nil {
    return nil, err
  }

  resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?appid=" + config.OpenWeatherMapApiKey + "&lat=" + lat + "&lon=" + long)
	if err != nil {
		log.Fatal("Could not get response", err)
		return nil, err
	}
	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Fatal("Could not decode body", err)
		return nil, err
	}

	return &d, nil
}

func weather(w http.ResponseWriter, r *http.Request) {
	city := r.PathValue("city")
	data, err := getWeather(city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/weather/{city}", weather)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
