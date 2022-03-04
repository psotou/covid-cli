package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

var (
	days   int
	region string
	comuna string
)

var (
	title  = color.New(color.FgWhite, color.Bold).Add(color.Underline)
	fields = color.New(color.FgWhite, color.Bold)
)

func main() {
	regionCmd := flag.NewFlagSet("region", flag.ExitOnError)
	comunaCmd := flag.NewFlagSet("comuna", flag.ExitOnError)

	regionCmd.StringVar(&region, "name", "metropolitana", "nombre de la región")
	regionCmd.IntVar(&days, "days", 7, "número de días")
	comunaCmd.StringVar(&comuna, "name", "nunoa", "nombre de la comuna")
	comunaCmd.IntVar(&days, "days", 7, "número de días")
	flag.Parse()

	args := flag.Args()
	switch args[0] {
	case "region":
		regionCmd.Parse(args[1:])
		casos, _ := CovidRegion(&days, &region)
		title.Printf("%s: %s\n", "Región", region)
		fields.Printf("%10s %9s %6s %6s\n", "Fecha", "Nacional", "Casos", "%")
		for i := range casos.Fechas {
			nacional, _ := strconv.ParseFloat(casos.Nacional[i], 64)
			region, _ := strconv.ParseFloat(casos.Region[i], 64)
			asPer := (region / nacional) * 100.0
			fmt.Printf("%10s %9.f %6s %6.1f\n", casos.Fechas[i], nacional, casos.Region[i], asPer)
		}
	case "comuna":
		comunaCmd.Parse(args[1:])
		casos, _ := CovidComuna(&days, &comuna)
		title.Printf("%s: %s\n", "Comuna", comuna)
		fields.Printf("%10s %6s\n", "Fecha", "Casos")
		// prints in reverse since the Comuna object return the data en in reverse order
		for i := range casos.Fechas {
			comunal, _ := strconv.ParseFloat(casos.Comuna[i], 64)
			fmt.Printf("%10s %6.f\n", casos.Fechas[i], comunal)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
