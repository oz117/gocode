package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var addr = flag.String("localhost", ":4000", "http service address") // http server address
var templ = template.Must(template.New("weather").Parse(TemplateStr))

/*
 * Retrieves the weatherData structure and shows the html page
 */
func getWeather(w http.ResponseWriter, r *http.Request) {
	mw := multiWeatherProvider{
		openWeatherMap{apikey: "your api key"},
		weatherUnderground{apikey: "your api key"},
	}
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := mw.temperature(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ.Execute(w, r.FormValue("s"))

	temp := strconv.FormatFloat(data, 'E', 2, 64)
	w.Write([]byte("<div>" + temp + "</div>"))
}

/*
 * Writes a simple welcome message
 */
func welcome(w http.ResponseWriter, r *http.Request) {
	templ.Execute(w, r.FormValue("s"))
	w.Write([]byte("<div>Welcome to my weather provider</div>"))
}

func main() {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/weather/", getWeather)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
