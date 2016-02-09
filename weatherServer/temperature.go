package main

import (
	"encoding/json"
	"log"
	"net/http"
)

/*
 * Simple interface
 */
type weatherProvider interface {
	temperature(city string) (float64, error)
}

/*
 * Structure for openweathermap.org
 */
type openWeatherMap struct {
	apikey string // api key necessary retrieve data from open weather map
}

/*
 * Structure for wunderground.com
 */
type weatherUnderground struct {
	apikey string
}

/*
 * A simple type that contains all the weather providers
 */
type multiWeatherProvider []weatherProvider

/*
 * Calls the temperature function of each weather provider
 * Returns the average temperature
 */
func (providers multiWeatherProvider) temperature(city string) (float64, error) {
	sum := 0.0

	for _, provider := range providers {
		k, err := provider.temperature(city)
		if err != nil {
			return 0, err
		}
		sum += k
	}
	return sum / float64(len(providers)), nil
}

/*
 * Retrieves the temperature from openweathermap.org
 * and returns it in celsius
 */
func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city + "&APPID=" + w.apikey)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}
	celsius := d.Main.Kelvin - 273.15
	log.Printf("openWeatherMap: %s: %.2f", city, celsius)
	return celsius, nil
}

/*
 * Retrieves the temperature from wunderground.com
 * and returns it in celsius
 */
func (w weatherUnderground) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.wunderground.com/api/" + w.apikey + "/conditions/q/" + city + ".json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var d struct {
		Observation struct {
			Celsius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}
	log.Printf("Weather Underground: %s: %2.f", city, d.Observation.Celsius)
	return d.Observation.Celsius, nil
}

/*
 * If we want to use pass each provider as a param and not as slice.
 * Works the same way as the temperature function taking a slice in parameter
 * Returns the average temperature in Celsius
 */
func temperature(city string, providers ...weatherProvider) (float64, error) {
	sum := 0.0

	for _, provider := range providers {
		c, err := provider.temperature(city)
		if err != nil {
			return 0, err
		}
		sum += c
	}
	return sum / float64(len(providers)), nil
}
