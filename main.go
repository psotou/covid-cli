package main

import (
	"covid-data/pkg/covid"
	"flag"
	"fmt"
)

var daysFlag *int

func main() {
	daysFlag = flag.Int("d", 7, "number of days to show")
	flag.Parse()

	casos := covid.Covid(daysFlag)

	// we pretty print the results jeje
	fmt.Println("-----------------------------------------------------------")
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", "Fecha", "Nacional", "RM", "Ñuñoa", "Provi", "Ñiquén", "Vallenar")
	fmt.Println("-----------------------------------------------------------")
	for i := range casos.Fecha {
		fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[i], casos.Nacional[i], casos.RM[i], casos.Nunoa[i], casos.Providencia[i], casos.Niquen[i], casos.Vallenar[i])
	}
	fmt.Println("-----------------------------------------------------------")
}
