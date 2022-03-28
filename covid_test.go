package main

import (
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestBaseData(t *testing.T) {
	type args struct {
		days int
	}
	tests := []struct {
		name string
		args args
		want *CasosCovid
	}{
		{
			name: "test",
			args: args{days: 7},
			want: mockCasosCovid(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseData(tt.args.days); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockCasosCovid() *CasosCovid {
	resp, err := http.Get(formatURL(nacional))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	d, _ := stringToLines(string(data))
	return &CasosCovid{
		Fechas:   lastValuesFromSlice(strings.Split(d[0], ","), 7),
		Nacional: strSlcToFloatSlc(lastValuesFromSlice(strings.Split(d[7], ","), 7)),
	}
}
