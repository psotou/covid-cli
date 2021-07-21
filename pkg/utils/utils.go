package utils

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func RetrieveData(url string) []byte {
	resp, err := http.Get(url)
	CheckErr("", err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	CheckErr("", err)

	return data
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
