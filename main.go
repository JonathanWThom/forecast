package main

import (
	"fmt"

	"github.com/jonathanwthom/forecast/internal"
)

func main() {
	periods := internal.GetForecast()
	for _, p := range periods {
		fmt.Println(p)
	}
}
