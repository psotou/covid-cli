package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	baseURL  = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/%s"
	nacional = "producto5/TotalesNacionales.csv"
	regional = "producto4/%s-CasosConfirmados-totalRegional.csv"
	comunal  = "producto1/Covid-19_std.csv"
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
		Fechas:   lastValuesSlice(strings.Split(data[0], ","), days),
		Nacional: strSlcToFloatSlc(lastValuesSlice(strings.Split(data[7], ","), days)),
	}
}

// AddDataRegional method adds upon BaseData object the corresponding number of cases according to a given region
func (cc *CasosCovid) AddsRegional(region *string) (CasosCovid, error) {
	for i := 1; i < len(cc.Fechas)+1; i++ {
		url := formatURL(fmt.Sprintf(regional, lastValue(cc.Fechas, i)))
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
			fecha := lastValue(cc.Fechas, i)
			if strings.Contains(v, fecha) && strings.Contains(v, strings.Title(*comuna)) {
				casosComuna.Fechas = append(casosComuna.Fechas, strings.Split(v, ",")[5])
				comuna, _ := strconv.ParseFloat(strings.Split(v, ",")[6], 64)
				casosComuna.Comuna = append(casosComuna.Comuna, comuna)
			}
		}
	}
	return casosComuna, nil
}

// TODO: refactor NacionalRegional and Comunal to printer functions that display the output
// I desire, which is the current text I output in main
func NacionalRegional(days *int, region *string) (CasosCovid, error) {
	nacionalRegional, err := BaseData(*days).AddsRegional(region)
	if err != nil {
		return CasosCovid{}, err
	}
	return nacionalRegional, nil
}

// TODO: I'm not much of fun of allocating a new CasosComuna struct just to reverse the print order.
// It should be a better way...
func Comunal(days *int, comuna *string) (CasosComuna, error) {
	comunal, err := BaseData(*days).DataComunal(comuna)
	if err != nil {
		return CasosComuna{}, err
	}
	casos := CasosComuna{}

	// inverse looping to return the last date first
	for i := len(comunal.Comuna) - 1; i >= 0; i-- {
		casos.Fechas = append(casos.Fechas, comunal.Fechas[i])
		casos.Comuna = append(casos.Comuna, comunal.Comuna[i])
	}
	return casos, nil
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

// lastValue returns the (n - pos) value of a given slice starting from the end
func lastValue(data []string, pos int) string {
	return data[len(data)-pos]
}

func lastValuesSlice(data []string, values int) []string {
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
