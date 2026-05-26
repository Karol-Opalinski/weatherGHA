// program wykonywalny
package main

// dołączenie bibliotek
import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// definicje struktur projektu
type City struct {
	Name      string  `json:"Name"`
	Country   string  `json:"Country"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type CitySelectModel struct {
	Name    string
	Country string
}

type WeatherInfoModel struct {
	LocalTime     string
	Temperature   float64
	WindSpeed     float64
	WindDirection float64
	Day           float64
	WeatherCode   float64
}

// tablica dostępnych miast (wczytywane z cities.json)
var cities []City

func loadCities() {

	data, err := os.ReadFile("cities.json")

	if err != nil {
		log.Fatal(err)
	}

	// mapowanie JSON na tablicę obiektów
	err = json.Unmarshal(data, &cities)

	if err != nil {
		log.Fatal(err)
	}

}

// handler dla wczytywania strony głównej root
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	countryMap := make(map[string]bool)
	var countries []string

	// wybranie unikalnych krajów
	for _, city := range cities {
		if !countryMap[city.Country] {
			countryMap[city.Country] = true
			countries = append(countries, city.Country)
		}
	}

	// stworzenie i wysłanie strony na bazie szablonu z przesłaną listą krajów dla Select {{range .}} -> iteracja po tablicy krajów
	// template.Must to odpowiednik 			tmpl, err := template.ParseFiles("index.html")
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, countries)
}

// handler dla przesłania informacji pogodowych dla miasta
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("city")
	if cityName == "" {
		http.Error(w, "Nie wybrano miasta", 400)
		return
	}

	for _, city := range cities {
		// pobranie informacji pogodowych dla miasta z API open-meteo
		if city.Name == cityName {
			foundCity := city
			url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true&timezone=auto", foundCity.Latitude, foundCity.Longitude)
			resp, err := http.Get(url)
			if err != nil {
				http.Error(w, "Błąd pobierania danych pogodowych", 500)
				return
			}

			// automatyczne zwolnienie zasobów odpowiedzi HTTP po zakończeniu funkcji
			defer resp.Body.Close()

			// zamiana JSON na mapę string:any
			var data map[string]any
			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				http.Error(w, "Błąd parsowania JSON", 500)
				return
			}

			// zamiana części z data na mapę string:any
			current := data["current_weather"].(map[string]any)

			// wybranie konretnych informacji pogodowych do przesłania
			weatherData := WeatherInfoModel{
				LocalTime:     current["time"].(string),
				Temperature:   current["temperature"].(float64),
				WindSpeed:     current["windspeed"].(float64),
				WindDirection: current["winddirection"].(float64),
				Day:           current["is_day"].(float64),
				WeatherCode:   current["weathercode"].(float64),
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(weatherData) // 200 ok
			return
		}
	}

	http.Error(w, "Nie znaleziono miasta", 404)

}

// handler do przesłania listy miast JSON dla wybranego kraju w celu uzupełnienia Select
func cityHandler(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		http.Error(w, "Nie wybrano kraju", 400)
		return
	}

	// pobranie odpowiednich miast dla kraju
	var result []CitySelectModel
	for _, city := range cities {
		if city.Country == country {
			cityData := CitySelectModel{
				Name:    city.Name,
				Country: city.Country,
			}
			result = append(result, cityData)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result) // 200 ok

}

func main() {
	// tryb healthcheck uruchamiany przez Docker HEALTHCHECK
	// jeżeli program zostanie uruchomiony z argumentem "--healthcheck"
	// proces kończy się natychmiast kodem 0, co oznacza poprawny stan aplikacji
	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		os.Exit(0)
	}

	// wczytanie miast z pliku cities.json
	loadCities()

	// udostępnienie folderów projektu
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	// podłączenie handlerów (router)
	http.HandleFunc("/", homePageHandler)
	http.HandleFunc("/cities", cityHandler)
	http.HandleFunc("/weather", weatherHandler)

	// log informacji dla konsoli
	log.SetFlags(0)
	log.Println("===================================")
	log.Println("Data uruchomienia:", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("Autor: Karol Opaliński")
	log.Println("Nasłuchiwanie na porcie TCP: 8080")
	log.Println("===================================")

	http.ListenAndServe(":8080", nil) // start serwera na porcie 8080
}
