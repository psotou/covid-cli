package datosmin

import (
	"covid-data/pkg/utils"
	"strings"
)

func Comuna(comuna string, lookingUpDates string, dataComunal []byte) string {
	grepDateBytes := utils.Grep(lookingUpDates, dataComunal)
	grepComunaBytes := utils.Grep(comuna, grepDateBytes)
	grepComunaResult := string(grepComunaBytes)
	totalCasesComuna := strings.Split(grepComunaResult, ",")

	// var casosComuna string
	if len(grepComunaBytes) == 0 {
		return "-"
	} else {
		return totalCasesComuna[6][:len(totalCasesComuna[6])-3]
	}
}
