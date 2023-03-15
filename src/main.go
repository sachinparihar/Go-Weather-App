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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		apiKey := "db9058e7659fb6c7faf51f7523524e45"
		city := r.URL.Query().Get("city")

		if city == "" {
			fmt.Fprint(w, "Please provide a city name")
			return
		}

		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, apiKey)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(w, "Weather data not found for city: %s", city)
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
		w.Write(jsonData)
	})

	log.Println("Server running on port 8000")
	http.ListenAndServe(":8000", nil)
}
