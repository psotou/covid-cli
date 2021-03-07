package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	baseURL  string = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output"
	nacional string = baseURL + "/producto5/TotalesNacionales.csv"
)

// CasosCovid is an object to save the results of the weeks requests
type CasosCovid struct {
	Fecha    []string
	Nacional []string
	RM       []string
}

func main() {
	resp, err := http.Get(nacional)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	dataStr := string(data)
	lines, err := StringToLines(dataStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	lastWeekDates := strings.Split(lines[0], ",")
	lastweekCases := strings.Split(lines[7], ",")

	// initializes the struct
	var casos *CasosCovid = new(CasosCovid)

	for i := 1; i < 8; i++ {
		// NACIONAL
		nationalDates := lastWeekDates[len(lastWeekDates)-i]
		nationalCases := lastweekCases[len(lastweekCases)-i][0:4]

		// REGIONAL
		fecha := nationalDates
		regional := baseURL + "/producto4/" + fecha + "-CasosConfirmados-totalRegional.csv"

		respRegional, err := http.Get(regional)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer respRegional.Body.Close()

		dataRegional, err := io.ReadAll(respRegional.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		dataStrRegional := string(dataRegional)
		linesRegional, err := StringToLines(dataStrRegional)
		if err != nil {
			log.Fatal(err.Error())
		}
		regionalCases := strings.Split(linesRegional[7], ",")

		// we append the results in the object casos
		casos.Fecha = append(casos.Fecha, nationalDates)
		casos.Nacional = append(casos.Nacional, nationalCases)
		casos.RM = append(casos.RM, regionalCases[9])

	}

	// we pretty print the results jeje
	fmt.Printf("%10s %10s %6s\n", "Fecha", "Nacional", "RM")
	fmt.Printf("%10s %10s %6s\n", casos.Fecha[0], casos.Nacional[0], casos.RM[0])
	fmt.Printf("%10s %10s %6s\n", casos.Fecha[1], casos.Nacional[1], casos.RM[1])
	fmt.Printf("%10s %10s %6s\n", casos.Fecha[2], casos.Nacional[2], casos.RM[2])
	fmt.Printf("%10s %10s %6s\n", casos.Fecha[3], casos.Nacional[3], casos.RM[3])
	fmt.Printf("%10s %10s %6s\n", casos.Fecha[4], casos.Nacional[4], casos.RM[4])
	fmt.Printf("%10s %10s %6s\n", casos.Fecha[5], casos.Nacional[5], casos.RM[5])
	fmt.Printf("%10s %10s %6s\n", casos.Fecha[6], casos.Nacional[6], casos.RM[6])
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
