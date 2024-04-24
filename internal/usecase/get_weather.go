package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tiagocosta/cloud-run/configs"
	"github.com/tiagocosta/cloud-run/internal/entity"
)

var cfg configs.Conf

type GetWeatherInputDTO struct {
	ZipCode string `json:"zipcode"`
}

type GetWeatherOutputDTO struct {
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

type GetWeatherUseCase struct {
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
}

type WeatherAPI struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64  `json:"temp_c"`
		Condition struct{} `json:"condition"`
	} `json:"current"`
}

func NewGetWeatherUseCase() *GetWeatherUseCase {
	cfg = configs.LoadConfig[configs.Conf](".")
	return &GetWeatherUseCase{}
}

func (uc *GetWeatherUseCase) Execute(input GetWeatherInputDTO) (GetWeatherOutputDTO, error) {
	weather, err := entity.NewWeather(input.ZipCode)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}
	viaCEP, err := FindZipCode(weather.Zip)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}

	weather.City = viaCEP.Localidade
	celsius, err := FetchWeather(weather.City)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}
	weather.FromCelsius(celsius)

	out := GetWeatherOutputDTO{
		Celsius:    fmt.Sprintf("%.1f", weather.Celsius),
		Fahrenheit: fmt.Sprintf("%.1f", weather.Fahrenheit),
		Kelvin:     fmt.Sprintf("%.1f", weather.Kelvin),
	}

	return out, nil
}

func FindZipCode(zipCode string) (*ViaCEP, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + zipCode + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c ViaCEP
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func FetchWeather(city string) (float64, error) {
	url := "https://api.weatherapi.com/v1/current.json?q=" + url.PathEscape(city) + "&key=" + cfg.WeatherAPIKey
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var w WeatherAPI
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return 0, err
	}

	if (WeatherAPI{}) == w {
		return 0, fmt.Errorf("can not find zipcode")
	}

	return w.Current.TempC, nil
}
