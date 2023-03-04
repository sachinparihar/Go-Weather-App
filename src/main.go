package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type weatherData struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temp"`
	Humidity    int     `json:"humidity"`
	Clouds      int     `json:"clouds"`
	WindSpeed   float64 `json:"wind_speed"`
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/weather/", handleWeather)
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(html))
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	// Extract city name from query parameter
	city := r.URL.Query().Get("city")
	if city == "" {
		fmt.Fprint(w, "Please provide a city name")
		return
	}

	// Construct URL to request weather data
	apiKey := "db9058e7659fb6c7faf51f7523524e45"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?APPID=%s&q=%s", apiKey, city)

	// Make request to OpenWeatherMap API
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Error fetching weather data: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading weather data: %v", err)
		return
	}

	// Parse JSON response into weatherData struct
	var data struct {
		Name string `json:"name"`
		Main struct {
			Kelvin   float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Fprintf(w, "Error decoding weather data: %v", err)
		return
	}

	// Set city name and return response as JSON
	weatherData := weatherData{
		City:        data.Name,
		Temperature: data.Main.Kelvin - 273.15,
		Humidity:    data.Main.Humidity,
		Clouds:      data.Clouds.All,
		WindSpeed:   data.Wind.Speed,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}
