package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func FormatURL(url string) string {
	return fmt.Sprintf(baseURL, url)
}

func RetrieveData(url string) ([]byte, error) {
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

func StringToLines(s string) ([]string, error) {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	return lines, err
}

func Grep(pattern string, data []byte) ([]byte, error) {
	grep := exec.Command("grep", pattern)
	grepIn, err := grep.StdinPipe()
	if err != nil {
		return nil, err
	}

	grepOut, err := grep.StdoutPipe()
	if err != nil {
		return nil, err
	}

	grep.Start()
	grepIn.Write(data)
	grepIn.Close()
	grepBytes, err := io.ReadAll(grepOut)
	if err != nil {
		return nil, err
	}
	grep.Wait()

	return grepBytes, nil
}
