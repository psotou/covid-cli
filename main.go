package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	days   *int
	region *string
	comuna *string
)

func main() {
	days = flag.Int("d", 7, "number of days to show")
	region = flag.String("r", "metropolitana", "region to show")
	comuna = flag.String("c", "nunoa", "comuna to show")
	flag.Parse()

	casos, err := Covid(days, region, comuna)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%s %s %s\n", "Fecha", "Nacional", *region)
	for i := range casos.Fechas {
		fmt.Printf("%s %s %s\n", casos.Fechas[i], casos.Nacional[i], casos.Region[i])
	}
}
