package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type WeatherData struct {
	Main struct {
		Temp     float32 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
}

func main() {
	listenAddr := ":8000"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/weather", HandleRequest)

	log.Printf("Starting server on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	apiKey := "db9058e7659fb6c7faf51f7523524e45"
	city := r.URL.Query().Get("city")

	if city == "" {
		http.Error(w, "Please provide a city name", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Weather data not found for city: %s", city), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"temp":       weatherData.Main.Temp,
		"humidity":   weatherData.Main.Humidity,
		"clouds":     weatherData.Clouds.All,
		"wind_speed": weatherData.Wind.Speed,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
