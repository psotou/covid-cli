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

	var casosComuna string
	if len(grepComunaBytes) == 0 {
		casosComuna = "-"
	} else {
		casosComuna = totalCasesComuna[1][0 : len(totalCasesComuna[1])-2]
	}
	return casosComuna
}
