package usecase

import (
	"fmt"

	"github.com/tiagocosta/cloud-run/configs"
	"github.com/tiagocosta/cloud-run/internal/entity"
	"github.com/tiagocosta/cloud-run/internal/pkg/wheatherapi"
	"github.com/tiagocosta/cloud-run/internal/pkg/zipcodeapi"
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

func NewGetWeatherUseCase() *GetWeatherUseCase {
	cfg = configs.LoadConfig[configs.Conf](".")
	return &GetWeatherUseCase{}
}

func (uc *GetWeatherUseCase) Execute(input GetWeatherInputDTO) (GetWeatherOutputDTO, error) {
	weather, err := entity.NewWeather(input.ZipCode)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}
	city, err := zipcodeapi.FindCity(weather.Zip)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}

	weather.City = city
	celsius, err := wheatherapi.GetWeather(weather.City, cfg.WeatherAPIKey)
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
