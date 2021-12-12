package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	baseURL  = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/%s"
	nacional = "producto5/TotalesNacionales.csv"
	regional = "producto4/%s-CasosConfirmados-totalRegional.csv"
	comunal  = "producto1/Covid-19_std.csv"
)

// CasosCovid is an object to save the results of the weeks requests
type CasosCovid struct {
	Fechas   []string
	Nacional []string
	Region   []string
	Comuna   []string
}

type Casos interface {
	DataNacional(*int) (CasosCovid, error)
	DataRegional(*string, *int) (CasosCovid, error)
	DataComunal(*string, *int) (CasosCovid, error)
}

// maps the regions with the index position in the object
// retrieved by the call to the region URL
var regiones = map[string]int{
	"arica":         1,
	"tarapaca":      2,
	"antofagaste":   3,
	"atacama":       4,
	"coquimbo":      5,
	"valparaiso":    6,
	"metropolitana": 7,
	"ohiggins":      8,
	"maule":         9,
	"nuble":         10,
	"biobio":        11,
	"araucania":     12,
	"losrios":       13,
	"loslagos":      14,
	"aysen":         15,
	"magallanes":    16,
}

func CovidData(url string) ([]string, error) {
	data, err := RetrieveData(url)
	if err != nil {
		return nil, err
	}
	return StringToLines(string(data))
}

// GetData retrieves all the data related to dates and new covid cases nation-wide
func GetData() Casos {
	data, err := CovidData(FormatURL(nacional))
	if err != nil {
		log.Fatalf(err.Error())
	}

	return &CasosCovid{
		Fechas:   strings.Split(data[0], ","),
		Nacional: strings.Split(data[7], ","),
	}
}

func (cc *CasosCovid) DataNacional(days *int) (CasosCovid, error) {
	// since I'm using Fechas and Nacional as a reference to work on the rest of the field
	// of the CasosCovid object, I cannot directly replace Fecha and Nacional fields
	// as done in the following DataNacional and DataComunal methods
	fechasRange := []string{}
	casosRange := []string{}
	for i := 1; i < *days+1; i++ {
		fechasRange = append(fechasRange, LastValue(cc.Fechas, i))
		casosRange = append(casosRange, LastValue(cc.Nacional, i))
	}
	return CasosCovid{
		Fechas:   fechasRange,
		Nacional: casosRange,
	}, nil
}

func (cc *CasosCovid) DataRegional(region *string, days *int) (CasosCovid, error) {
	for i := 1; i < *days+1; i++ {
		url := FormatURL(fmt.Sprintf(regional, LastValue(cc.Fechas, i)))
		data, err := CovidData(url)
		if err != nil {
			return CasosCovid{}, fmt.Errorf("error retrieving regional data for the date %s", LastValue(cc.Fechas, i))
		}
		// r is the number associated with the region
		r := regiones[*region]
		// the 9th position return the values for Casos Nuevos Totales
		cc.Region = append(cc.Region, strings.Split(data[r], ",")[9])
	}
	return *cc, nil
}

func (cc *CasosCovid) DataComunal(comuna *string, days *int) (CasosCovid, error) {
	data, err := CovidData(FormatURL(comunal))
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, v := range data {
		for i := 1; i < *days+1; i++ {
			fecha := LastValue(cc.Fechas, i)
			if strings.Contains(v, fecha) && strings.Contains(v, strings.Title(*comuna)) {
				cc.Comuna = append(cc.Comuna, v)
			}
		}
	}
	return *cc, nil
}

func CovidRegion(days *int, region *string) (CasosCovid, error) {
	nacional, _ := GetData().DataNacional(days)
	regional, _ := GetData().DataRegional(region, days)

	return CasosCovid{
		Fechas:   nacional.Fechas,
		Nacional: nacional.Nacional,
		Region:   regional.Region,
	}, nil
}

func CovidComuna(days *int, comuna *string) (CasosCovid, error) {
	comunal, _ := GetData().DataComunal(comuna, days)
	casos := CasosCovid{}
	for _, v := range comunal.Comuna {
		casos.Fechas = append(casos.Fechas, strings.Split(v, ",")[5])
		casos.Comuna = append(casos.Comuna, strings.Split(v, ",")[6])
	}
	return casos, nil
}
