package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
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
	return lines, scanner.Err()
}

// LastValue returns the (n - pos) value of a given slice
func LastValue(data []string, pos int) string {
	return data[len(data)-pos]
}
