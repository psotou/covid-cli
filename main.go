package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	days   int
	region string
	comuna string
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
		casos, _ := NacionalRegional(&days, &region)

		title.Printf("%s: %s\n", "Región", region)
		fields.Printf("%10s %9s %6s %6s\n", "Fecha", "Nacional", "Casos", "%")
		for i := range casos.Fechas {
			asPer := (casos.Region[i] / casos.Nacional[i]) * 100.0
			fmt.Printf("%10s %9.f %6.f %6.1f\n", casos.Fechas[i], casos.Nacional[i], casos.Region[i], asPer)
		}
	case "comuna":
		comunaCmd.Parse(args[1:])
		casos, _ := Comunal(&days, &comuna)

		title.Printf("%s: %s\n", "Comuna", comuna)
		fields.Printf("%10s %6s %4s\n", "Fecha", "Casos", "Delta")
		for i := range casos.Fechas {
			if i+1 < len(casos.Comuna) && len(casos.Comuna) > 1 && casos.Comuna[i] != 0 {
				fmt.Printf("%10s %6.f %5.f\n", casos.Fechas[i], casos.Comuna[i], casos.Comuna[i]-casos.Comuna[i+1])
			} else {
				fmt.Printf("%10s %6.f %5s\n", casos.Fechas[i], casos.Comuna[i], "--")
			}
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
