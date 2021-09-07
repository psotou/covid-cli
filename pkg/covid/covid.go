package covid

import (
	"covid-data/pkg/datosmin"
	"covid-data/pkg/utils"
	"fmt"
	"strconv"
	"strings"
)

const (
	baseURL    = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output"
	nacional   = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/producto5/TotalesNacionales.csv"
	regional   = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/producto4/%s-CasosConfirmados-totalRegional.csv"
	comunal    = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/producto6/bulk/data.csv"
	comunalDos = "https://raw.githubusercontent.com/MinCiencia/Datos-COVID19/master/output/producto1/Covid-19_std.csv"
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

func Covid(daysFlag *int) *CasosCovid {
	dataStr := string(utils.RetrieveData(nacional))
	lines, err := utils.StringToLines(dataStr)
	utils.CheckErr("", err)

	nacionalFechas := strings.Split(lines[0], ",")
	totalCases := strings.Split(lines[7], ",")

	// initializes the struct
	var casos *CasosCovid = new(CasosCovid)

	// we go from 1 to 7 since we only want data for just a week
	for i := 1; i < *daysFlag+1; i++ {
		// NACIONAL
		nationalDates := nacionalFechas[len(nacionalFechas)-i]
		nationalCases := totalCases[len(totalCases)-i]
		nationalCasesFloat, _ := strconv.ParseFloat(nationalCases, 64)

		// REGIONAL
		regional := fmt.Sprintf(regional, nationalDates)
		dataStrRegional := string(utils.RetrieveData(regional))
		linesRegional, err := utils.StringToLines(dataStrRegional)
		utils.CheckErr("", err)

		regionalCases := strings.Split(linesRegional[7], ",")

		casos.Fecha = append(casos.Fecha, nationalDates)
		casos.Nacional = append(casos.Nacional, fmt.Sprintf("%.f", nationalCasesFloat))
		casos.RM = append(casos.RM, regionalCases[9])

		// COMUNAL ACUMULADO
		dataComunal := utils.RetrieveData(comunalDos)
		// nationalDates := strings.Replace(nationalDates, "-", "/", -1)

		casos.Nunoa = append(casos.Nunoa, datosmin.Comuna("Nunoa", nationalDates, dataComunal))
		casos.Providencia = append(casos.Providencia, datosmin.Comuna("Providencia", nationalDates, dataComunal))
		casos.Niquen = append(casos.Niquen, datosmin.Comuna("Niquen", nationalDates, dataComunal))
		casos.Vallenar = append(casos.Vallenar, datosmin.Comuna("Vallenar", nationalDates, dataComunal))
	}
	return casos
}
