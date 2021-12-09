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

// maps the regions with the v position in the object
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

// GetNacional retrieves all the data related to dates and new covid cases nation-wide
func GetData() Casos {
	url := FormatURL(nacional)
	data, err := CovidData(url)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return &CasosCovid{
		Fechas:   strings.Split(data[0], ","),
		Nacional: strings.Split(data[7], ","),
	}
}

func (cc *CasosCovid) DataNacional(days *int) (CasosCovid, error) {
	fechasRange := []string{}
	casosRange := []string{}
	for i := 1; i < *days+1; i++ {
		fechasRange = append(fechasRange, cc.Fechas[len(cc.Fechas)-i])
		casosRange = append(casosRange, cc.Nacional[len(cc.Nacional)-i])
	}

	return CasosCovid{
		Fechas:   fechasRange,
		Nacional: casosRange,
	}, nil
}

func (cc *CasosCovid) DataRegional(region *string, days *int) (CasosCovid, error) {
	casosRegional := []string{}
	for i := 1; i < *days+1; i++ {
		url := FormatURL(fmt.Sprintf(regional, cc.Fechas[len(cc.Fechas)-i]))
		dataRegional, err := CovidData(url)
		if err != nil {
			return CasosCovid{}, fmt.Errorf("error retrieving regional data for the date %s", cc.Fechas[len(cc.Fechas)-i])
		}
		// r is the number associated to the region
		r := regiones[*region]
		// the 9th position return the values for Casos Nuevos Totales
		casosRegional = append(casosRegional, strings.Split(dataRegional[r], ",")[9])
	}

	return CasosCovid{
		Region: casosRegional,
	}, nil
}

func (cc *CasosCovid) DataComunal(comuna *string, days *int) (CasosCovid, error) {
	url := FormatURL(comunal)
	data, err := CovidData(url)
	if err != nil {
		log.Fatalf(err.Error())
	}
	casosComuna := []string{}
	for _, v := range data {
		for i := 1; i < *days+1; i++ {
			fecha := cc.Fechas[len(cc.Fechas)-i]
			if strings.Contains(v, fecha) && strings.Contains(v, strings.Title(*comuna)) {
				// casosComuna = append(casosComuna, strings.Split(v, ","))
				casosComuna = append(casosComuna, v)
			}
		}
	}
	return CasosCovid{
		Comuna: casosComuna,
	}, nil
}

func Covid(days *int, region, comuna *string) (CasosCovid, error) {
	nacional, _ := GetData().DataNacional(days)
	regional, _ := GetData().DataRegional(region, days)
	comunal, _ := GetData().DataComunal(comuna, days)

	return CasosCovid{
		Fechas:   nacional.Fechas,
		Nacional: nacional.Nacional,
		Region:   regional.Region,
		Comuna:   comunal.Comuna,
	}, nil
}
