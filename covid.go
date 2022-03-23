package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	baseURL  = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/%s"
	nacional = "producto5/TotalesNacionales.csv"
	regional = "producto4/%s-CasosConfirmados-totalRegional.csv"
	comunal  = "producto1/Covid-19_std.csv"
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
	data, err := CovidData(formatURL(nacional))
	if err != nil {
		log.Fatalf(err.Error())
	}
	return &CasosCovid{
		Fechas:   lastValuesFromSlice(strings.Split(data[0], ","), days),
		Nacional: strSlcToFloatSlc(lastValuesFromSlice(strings.Split(data[7], ","), days)),
	}
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
		// the 8th position return the values for Casos Nuevos Totales
		// which is the sum of the Casos Nuevos con Síntomas, Casos Nuevos sin Síntomas, Casos Nuevos Reportados por Laboratorio
		region, _ := strconv.ParseFloat(strings.Split(data[r], ",")[8], 64)
		cc.Region = append(cc.Region, region)
	}
	return *cc, nil
}

func (cc *CasosCovid) DataComunal(comuna *string) (CasosComuna, error) {
	data, err := CovidData(formatURL(comunal))
	if err != nil {
		return CasosComuna{}, err
	}
	casosComuna := CasosComuna{}
	for _, v := range data {
		for i := 1; i < len(cc.Fechas)+1; i++ {
			fecha := valueAtNthPosFromEnd(cc.Fechas, i)
			if strings.Contains(v, fecha) && strings.Contains(v, strings.Title(*comuna)) {
				comuna, _ := strconv.ParseFloat(strings.Split(v, ",")[6], 64)

				casosComuna.Fechas = append(casosComuna.Fechas, strings.Split(v, ",")[5])
				casosComuna.Comuna = append(casosComuna.Comuna, comuna)
			}
		}
	}
	return casosComuna, nil
}

type Print struct{}

func (p *Print) DataNacionalRegional(casos CasosCovid, days *int, region *string) {
	title.Printf("%s: %s\n", "Región", *region)
	fields.Printf("%10s %9s %6s %6s\n", "Fecha", "Nacional", "Casos", "%")

	for i := len(casos.Fechas) - 1; i >= 0; i-- {
		regToNacPer := (casos.Region[i] / casos.Nacional[i]) * 100.0
		fmt.Printf("%10s %9.f %6.f %6.1f\n", casos.Fechas[i], casos.Nacional[i], casos.Region[i], regToNacPer)
	}
}

func (p *Print) DataComunal(casos CasosComuna, days *int, comuna *string) {
	title.Printf("%s: %s\n", "Comuna", *comuna)
	fields.Printf("%10s %6s %4s\n", "Fecha", "Casos", "Delta")

	for i := len(casos.Fechas) - 1; i >= 0; i-- {
		if len(casos.Comuna) > 0 && i > 0 {
			delta := casos.Comuna[i] - casos.Comuna[i-1]
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
