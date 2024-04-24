package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

type WeatherAPIMocked struct {
	mock.Mock
}

func (m *WeatherAPIMocked) FetchWeather(city string) (float64, error) {
	args := m.Called(city)
	return args.Get(0).(float64), args.Error(1)
}

func TestFindZipCode(t *testing.T) {
	type args struct {
		zipCode string
	}
	tests := []struct {
		name    string
		args    args
		want    *ViaCEP
		wantErr bool
	}{
		{
			"Invalid zipcode (length less than 8). Should return an error",
			args{
				zipCode: "1234567",
			},
			nil,
			true,
		},
		{
			"Valid but non existing zipcode. Should return an empty struct",
			args{
				zipCode: "12345678",
			},
			&ViaCEP{},
			false,
		},
		{
			"Valid and existing zipcode. Should return an fully populated struct",
			args{
				zipCode: "71218010",
			},
			&ViaCEP{
				Cep:         "71218-010",
				Logradouro:  "SMAS Trecho 1 C",
				Complemento: "(Condomínio Living SQPS)",
				Bairro:      "Zona Industrial (Guará)",
				Localidade:  "Brasília",
				Uf:          "DF",
				Ibge:        "5300108",
				Gia:         "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindZipCode(tt.args.zipCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindZipCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindZipCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchWeather(t *testing.T) {
	mock := new(WeatherAPIMocked)
	mock.On("FetchWeather", "Brasília").Return(29.0, nil)
	mock.On("FetchWeather", "non existing city").Return(0.0, errors.New("can not find zipcode"))
	type args struct {
		city string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			"Valid City. Should return the correct temperature",
			args{city: "Brasília"},
			29.0,
			false,
		},
		{
			"Invalid City. Should return error",
			args{city: "non existing city"},
			0.0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mock.FetchWeather(tt.args.city)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWeather() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FetchWeather() = %v, want %v", got, tt.want)
			}
		})
	}
}
