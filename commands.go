package main

import (
	"flag"
	"fmt"
)

func Run() error {
	regionFlagSet := flag.NewFlagSet("region", flag.ExitOnError)
	comunaFlagSet := flag.NewFlagSet("comuna", flag.ExitOnError)

	regionName := regionFlagSet.String("name", "metropolitana", "-name [nombre región]")
	regionDays := regionFlagSet.Int("days", 7, "-days [número de días]")
	comunaName := comunaFlagSet.String("name", "nunoa", "-name [nombre comuna]")
	comunaDays := comunaFlagSet.Int("days", 7, "-days [número de días]")
	flag.Parse()

	print := Print{}
	var cmd string
	var args []string
	if len(flag.Args()) > 0 {
		cmd, args = flag.Args()[0], flag.Args()[1:]
	}

	switch cmd {
	case "region":
		regionFlagSet.Parse(args)
		casos, err := BaseData(*regionDays).AddDataRegional(regionName)
		if err != nil {
			return err
		}
		print.DataNacionalRegional(casos, regionName)
	case "comuna":
		comunaFlagSet.Parse(args)
		casos, err := BaseData(*comunaDays).DataComunal(comunaName)
		if err != nil {
			return err
		}
		print.DataComunal(casos, comunaName)
	default:
		if cmd == "" || cmd == "help" {
			usage()
			return flag.ErrHelp
		}
		return fmt.Errorf("covid %s: unkown option", cmd)
	}
	return nil
}

func usage() {
	fmt.Println(`
covid es una herramienta para consultar los casos diarios de covid19 a nivel nacional, regional y comunal.

Uso:

    covid <option> [arguments]

option:
    
    region    para consultar el número de casos diaros por región y contrastarlos con el toral diario nacional
    comuna    para consultar el número de casos diarios por comuna. Muestra el acumulado a la fecha.

arguments:

    days      para consultar un número determinado de días. Por defecto consulta 7 días.
    name      para consultar una región o comuna determinada. La región por defecto es la RM, y la comuna Ñuñoa.
        `)
}
