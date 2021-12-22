package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LatLong struct {
	Lat float64
	Lon float64
}

type Period struct {
	Name             string
	Temperature      int
	TemperatureUnit  string
	ShortForecast    string
	DetailedForecast string
}

func (p Period) String() string {
	return fmt.Sprintf("%s: %d%s. %s.", p.Name, p.Temperature, p.TemperatureUnit, p.ShortForecast)
}

type Properties struct {
	Forecast       string
	ForecastHourly string
	Periods        []Period
}

type Points struct {
	Properties Properties
}

func GetForecast() []Period {
	var data Points
	resp, err := http.Get(GetQuadrant())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(resp.Body, &buf))
	if err = dec.Decode(&data); err != nil {
		log.Fatal(err)
	}

	return data.Properties.Periods
}

func GetQuadrant() string {
	base := "https://api.weather.gov/points/"
	latlong, err := GetLocationData()
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%v%s", base, latlong)
	var data Points
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(resp.Body, &buf))
	if err = dec.Decode(&data); err != nil {
		log.Fatal(err)
	}

	return data.Properties.Forecast
}

func (l LatLong) String() string {
	return fmt.Sprintf("%v,%v", l.Lat, l.Lon)
}

func GetLocationData() (LatLong, error) {
	var data LatLong
	url := "http://ip-api.com/json/"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(resp.Body, &buf))
	if err = dec.Decode(&data); err != nil {
		log.Fatal(err)
	}

	return data, nil
}
