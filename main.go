package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
	baseURL  string = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output"
	nacional string = baseURL + "/producto5/TotalesNacionales.csv"
	comunal  string = baseURL + "/producto6/bulk/data.csv"
)

// CasosCovid is an object to save the results of the weeks requests
type CasosCovid struct {
	Fecha       []string
	Nacional    []string
	RM          []string
	Nunoa       []string
	Providencia []string
	Niquen      []string
	Vallenar    []string
}

func main() {
	resp, err := http.Get(nacional)
	CheckErr("", err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	CheckErr("", err)

	dataStr := string(data)
	lines, err := StringToLines(dataStr)
	CheckErr("", err)

	nacionalFechas := strings.Split(lines[0], ",")
	totalCases := strings.Split(lines[7], ",")

	// initializes the struct
	var casos *CasosCovid = new(CasosCovid)

	for i := 1; i < 8; i++ {
		// NACIONAL
		nationalDates := nacionalFechas[len(nacionalFechas)-i]
		nationalCases := totalCases[len(totalCases)-i][0:4]

		// REGIONAL
		fecha := nationalDates
		regional := baseURL + "/producto4/" + fecha + "-CasosConfirmados-totalRegional.csv"

		respRegional, err := http.Get(regional)
		CheckErr("", err)
		defer respRegional.Body.Close()

		dataRegional, err := io.ReadAll(respRegional.Body)
		CheckErr("", err)

		dataStrRegional := string(dataRegional)
		linesRegional, err := StringToLines(dataStrRegional)
		CheckErr("", err)

		regionalCases := strings.Split(linesRegional[7], ",")

		// we append the results in the object casos
		casos.Fecha = append(casos.Fecha, nationalDates)
		casos.Nacional = append(casos.Nacional, nationalCases)
		casos.RM = append(casos.RM, regionalCases[9])

		// COMUNAL ACUMULADO
		respComunal, err := http.Get(comunal)
		CheckErr("", err)
		defer respComunal.Body.Close()

		dataComunal, err := io.ReadAll(respComunal.Body)
		CheckErr("", err)

		lookingUpDates := strings.Replace(nationalDates, "-", "/", -1)
		// grepDateBytes := Grep(lookingUpDates, dataComunal)

		// Ñuñoa
		casos.Nunoa = append(casos.Nunoa, Comuna("Ñuñoa", lookingUpDates, dataComunal))

		// Providencia
		casos.Providencia = append(casos.Providencia, Comuna("Providencia", lookingUpDates, dataComunal))

		// Ñiquén
		casos.Niquen = append(casos.Niquen, Comuna("Ñiquén", lookingUpDates, dataComunal))

		// Vallenar
		casos.Vallenar = append(casos.Vallenar, Comuna("Vallenar", lookingUpDates, dataComunal))

	}

	// we pretty print the results jeje
	fmt.Println("-----------------------------------------------------------")
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", "Fecha", "Nacional", "RM", "Ñuñoa", "Provi", "Ñiquén", "Vallenar")
	fmt.Println("-----------------------------------------------------------")
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[0], casos.Nacional[0], casos.RM[0], casos.Nunoa[0], casos.Providencia[0], casos.Niquen[0], casos.Vallenar[0])
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[1], casos.Nacional[1], casos.RM[1], casos.Nunoa[1], casos.Providencia[1], casos.Niquen[1], casos.Vallenar[1])
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[2], casos.Nacional[2], casos.RM[2], casos.Nunoa[2], casos.Providencia[2], casos.Niquen[2], casos.Vallenar[2])
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[3], casos.Nacional[3], casos.RM[3], casos.Nunoa[3], casos.Providencia[3], casos.Niquen[3], casos.Vallenar[3])
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[4], casos.Nacional[4], casos.RM[4], casos.Nunoa[4], casos.Providencia[4], casos.Niquen[4], casos.Vallenar[4])
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[5], casos.Nacional[5], casos.RM[5], casos.Nunoa[5], casos.Providencia[5], casos.Niquen[5], casos.Vallenar[5])
	fmt.Printf("%10s %10s %6s %6s %6s %6s %8s\n", casos.Fecha[6], casos.Nacional[6], casos.RM[6], casos.Nunoa[6], casos.Providencia[6], casos.Niquen[6], casos.Vallenar[6])
	fmt.Println("-----------------------------------------------------------")
}

// StringToLines reads lines from a string
func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func Comuna(comuna string, lookingUpDates string, dataComunal []byte) string {
	grepDateBytes := Grep(lookingUpDates, dataComunal)
	grepComunaBytes := Grep(comuna, grepDateBytes)
	grepComunaResult := string(grepComunaBytes)
	totalCasesComuna := strings.Split(grepComunaResult, ",")

	var casosComuna string
	if len(grepComunaBytes) == 0 {
		casosComuna = "-"
	} else {
		casosComuna = totalCasesComuna[1][0 : len(totalCasesComuna[1])-2]
	}
	return casosComuna
}

func Grep(pattern string, data []byte) []byte {
	grep := exec.Command("grep", pattern)
	grepIn, err := grep.StdinPipe()
	CheckErr("", err)

	grepOut, err := grep.StdoutPipe()
	CheckErr("", err)

	grep.Start()
	grepIn.Write(data)
	grepIn.Close()
	grepBytes, err := io.ReadAll(grepOut)
	CheckErr("", err)
	grep.Wait()

	return grepBytes
}

// CheckErr to handle error
func CheckErr(str string, err error) {
	if err != nil {
		log.Fatal(str, err.Error())
	}
}
