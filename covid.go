package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	baseURL     = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/%s"
	nacionalStd = "producto5/TotalesNacionales_std.csv"
	regional    = "producto4/%s-CasosConfirmados-totalRegional.csv"
	comunalStd  = "producto1/Covid-19_std.csv"
)

var (
	title  = color.New(color.FgWhite, color.Bold).Add(color.Underline)
	fields = color.New(color.FgWhite, color.Bold)
)

type CasosCovid struct {
	Fechas   []string
	Nacional []float64
	Region   []float64
}

type CasosComuna struct {
	Fechas []string
	Comuna []float64
}

// maps the regions with the index position in the object retrieved by the call to the region URL
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
	data, err := retrieveData(url)
	if err != nil {
		return nil, err
	}
	return stringToLines(string(data))
}

// BaseData works as a sort of constructor that initializes the CasosCovid struct and populates it
// with the Fechas and Nacional number of covid cases according to the given number of days
func BaseData(days int) *CasosCovid {
	data, err := CovidData(formatURL(nacionalStd))
	if err != nil {
		log.Fatalf(err.Error())
	}

	var casos CasosCovid
	for i := len(data) - 1; i >= 0; i-- {
		for j := 0; j < days; j++ {
			date := time.Now().AddDate(0, 0, -j).Format("2006-01-02")
			if strings.Contains(data[i], date) && strings.Contains(data[i], "Casos nuevos totales") {
				// at 2nd position we found the value for the amount of Casos Nuevos Totales
				casosInt, _ := strconv.ParseFloat(strings.Split(data[i], ",")[2], 64)
				casos.Nacional = append(casos.Nacional, casosInt)
				casos.Fechas = append(casos.Fechas, date)
			}
		}
	}

	return &casos
}

// AddDataRegional method adds upon BaseData object the corresponding number of cases according to a given region
func (cc *CasosCovid) AddDataRegional(region *string) (CasosCovid, error) {
	for i := 1; i < len(cc.Fechas)+1; i++ {
		url := formatURL(fmt.Sprintf(regional, valueAtNthPosFromEnd(cc.Fechas, i)))
		data, err := CovidData(url)
		if err != nil {
			return CasosCovid{}, err
		}
		// r is the number associated with the region
		r := regiones[*region]
		// the 10th position return the values for Casos Nuevos Totales
		// which is the sum of the Casos Nuevos con Síntomas, Casos Nuevos sin Síntomas, Casos Nuevos Reportados por Laboratorio
		region, _ := strconv.ParseFloat(strings.Split(data[r], ",")[10], 64)
		cc.Region = append(cc.Region, region)
	}
	return *cc, nil
}

func (cc *CasosCovid) DataComunal(comuna *string) (CasosComuna, error) {
	data, err := CovidData(formatURL(comunalStd))
	if err != nil {
		return CasosComuna{}, err
	}
	casosComuna := CasosComuna{}
	for i := len(data) - 1; i >= 0; i-- {
		for j := 1; j < len(cc.Fechas)+1; j++ {
			fecha := valueAtNthPosFromEnd(cc.Fechas, j)
			if strings.Contains(data[i], fecha) && strings.Contains(data[i], strings.Title(*comuna)) {
				// at 6th position we found the amount of cases per comuna
				comuna, _ := strconv.ParseFloat(strings.Split(data[i], ",")[6], 64)
				// at 5th position we found the dates
				casosComuna.Fechas = append(casosComuna.Fechas, strings.Split(data[i], ",")[5])
				casosComuna.Comuna = append(casosComuna.Comuna, comuna)
			}
		}
	}
	return casosComuna, nil
}

type Print struct{}

func (p *Print) DataNacionalRegional(casos CasosCovid, region *string) {
	title.Printf("%s: %s\n", "Región", *region)
	fields.Printf("%10s %9s %6s %6s\n", "Fecha", "Nacional", "Casos", "%")

	for i := range casos.Fechas {
		regToNacPer := (casos.Region[len(casos.Region)-1-i] / casos.Nacional[i]) * 100.0
		casosRegion := casos.Region[len(casos.Region)-1-i]
		fmt.Printf("%10s %9.f %6.f %6.1f\n", casos.Fechas[i], casos.Nacional[i], casosRegion, regToNacPer)
	}
}

func (p *Print) DataComunal(casos CasosComuna, comuna *string) {
	title.Printf("%s: %s\n", "Comuna", *comuna)
	fields.Printf("%10s %6s %4s\n", "Fecha", "Casos", "Delta")

	for i := range casos.Comuna {
		if i+1 < len(casos.Comuna) {
			delta := casos.Comuna[i] - casos.Comuna[i+1]
			fmt.Printf("%10s %6.f %5.f\n", casos.Fechas[i], casos.Comuna[i], delta)
		} else {
			fmt.Printf("%10s %6.f %5s\n", casos.Fechas[i], casos.Comuna[i], "--")
		}
	}
}

func formatURL(url string) string {
	return fmt.Sprintf(baseURL, url)
}

func retrieveData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func stringToLines(s string) ([]string, error) {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// returns the value of a slice at a position pos counting from end to start
func valueAtNthPosFromEnd(data []string, pos int) string {
	return data[len(data)-pos]
}

func lastValuesFromSlice(data []string, values int) []string {
	return data[len(data)-values:]
}

func strSlcToFloatSlc(slc []string) []float64 {
	fSlc := []float64{}
	for _, val := range slc {
		fVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fSlc = append(fSlc, fVal)
	}
	return fSlc
}
