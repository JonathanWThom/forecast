package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Forecast struct {
	City        string
	State       string
	ForecastUrl string
	Periods     []Period
}

func NewForecast() Forecast {
	f := Forecast{}
	f.SetQuadrantParameters()
	f.SetPeriods()

	return f
}

func (f Forecast) String() string {
	header := fmt.Sprintf(
		"Forecast for %s, %s",
		f.City, f.State,
	)
	var body string
	for _, p := range f.Periods {
		body += fmt.Sprintf("\n%s", p)
	}

	return fmt.Sprintf("%s\n%s", header, body)
}

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

type LocationProperties struct {
	City  string
	State string
}

type RelativeLocation struct {
	Properties LocationProperties
}

type Properties struct {
	Forecast         string
	ForecastHourly   string
	Periods          []Period
	RelativeLocation RelativeLocation
}

type Points struct {
	Properties Properties
}

func (f *Forecast) SetPeriods() {
	var data Points
	resp, err := http.Get(f.ForecastUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(resp.Body, &buf))
	if err = dec.Decode(&data); err != nil {
		log.Fatal(err)
	}

	f.Periods = data.Properties.Periods
}

func (f *Forecast) SetQuadrantParameters() {
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

	f.ForecastUrl = data.Properties.Forecast
	f.City = data.Properties.RelativeLocation.Properties.City
	f.State = data.Properties.RelativeLocation.Properties.State
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
