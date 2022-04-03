package main

import (
	"flag"
	"log"
	"os"
)

var (
	days   int
	region string
	comuna string
)

func main() {
	regionCmd := flag.NewFlagSet("region", flag.ExitOnError)
	comunaCmd := flag.NewFlagSet("comuna", flag.ExitOnError)

	regionCmd.StringVar(&region, "name", "metropolitana", "nombre de la región")
	regionCmd.IntVar(&days, "days", 7, "número de días")
	comunaCmd.StringVar(&comuna, "name", "nunoa", "nombre de la comuna")
	comunaCmd.IntVar(&days, "days", 7, "número de días")
	flag.Parse()

	print := Print{}

	args := flag.Args()
	switch args[0] {
	case "region":
		regionCmd.Parse(args[1:])
		casos, err := BaseData(days).AddDataRegional(&region)
		if err != nil {
			log.Fatalf(err.Error())
		}
		print.DataNacionalRegional(casos, &region)
	case "comuna":
		comunaCmd.Parse(args[1:])
		casos, err := BaseData(days).DataComunal(&comuna)
		if err != nil {
			log.Fatalf(err.Error())
		}
		print.DataComunal(casos, &comuna)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
