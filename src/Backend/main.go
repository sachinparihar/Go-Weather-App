package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	fileServer := http.FileServer(http.Dir("../Frontend/"))
	http.Handle("/", fileServer)
	http.HandleFunc("/weather", HandleWeather)
	listenAddr := ":8000"
	log.Printf("Starting server on %s", listenAddr)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func HandleWeather(w http.ResponseWriter, r *http.Request) {
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
		"city":       city,
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
	_, err = w.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
}
