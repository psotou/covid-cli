package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	days   int
	region string
	comuna string
)

func main() {
	regionCmd := flag.NewFlagSet("region", flag.ExitOnError)
	comunaCmd := flag.NewFlagSet("comuna", flag.ExitOnError)

	regionCmd.StringVar(&region, "name", "metropolitana", "region")
	regionCmd.IntVar(&days, "days", 7, "last given days")
	comunaCmd.StringVar(&comuna, "name", "nunoa", "region")
	comunaCmd.IntVar(&days, "days", 7, "last given days")
	flag.Parse()

	args := flag.Args()
	switch args[0] {
	case "region":
		regionCmd.Parse(args[1:])

		casos, _ := CovidRegion(&days, &region)
		fmt.Printf("%10s %10s %6s\n", "Fecha", "Nacional", strings.Title(region)[0:5])
		for i := range casos.Fechas {
			nacional, _ := strconv.ParseFloat(casos.Nacional[i], 64)
			fmt.Printf("%10s %10.f %6s\n", casos.Fechas[i], nacional, casos.Region[i])
		}
	case "comuna":
		comunaCmd.Parse(args[1:])

		casos, _ := CovidComuna(&days, &comuna)
		fmt.Printf("%10s %6s\n", "Fecha", strings.Title(comuna)[0:5])
		for i := range casos.Fechas {
			comunal, _ := strconv.ParseFloat(casos.Comuna[i], 64)
			fmt.Printf("%10s %6.f\n", casos.Fechas[i], comunal)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
