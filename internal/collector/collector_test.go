package collector

import (
	"reflect"
	"testing"
)



func Test_convertHeaders(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		want    []float64
		wantErr bool
	}{
		{"Normalize header",  "100;w=21600", []float64{100,21600},false},
		{"empty header",  "", nil,true},
	}

		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertHeaders(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertHeaders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToFloat(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []float64
		wantErr bool
	}{

		{"strings to float",[]string{"100","26000"},[]float64{100,26000},false},
		{"one item",  []string{"100"}, []float64{100},false},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToFloat(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
