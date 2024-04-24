package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tiagocosta/cloud-run/configs"
	"github.com/tiagocosta/cloud-run/internal/entity"
)

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
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func NewGetWeatherUseCase() *GetWeatherUseCase {
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
	err = FetchWeather(weather)
	if err != nil {
		return GetWeatherOutputDTO{}, err
	}

	out := GetWeatherOutputDTO{
		Celsius:    fmt.Sprintf("%.1f", weather.Celsius),
		Fahrenheit: fmt.Sprintf("%.1f", weather.Fahrenheit),
		Kelvin:     fmt.Sprintf("%.1f", weather.Kelvin),
	}

	return out, nil
}

func FindZipCode(zipCode string) (*ViaCEP, error) {
	fmt.Println("https://viacep.com.br/ws/" + zipCode + "/json/")
	resp, error := http.Get("https://viacep.com.br/ws/" + zipCode + "/json/")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	var c ViaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}

	return &c, nil
}

func FetchWeather(weather *entity.Weather) error {
	resp, error := http.Get("http://api.weatherapi.com/v1/current.json?key=" + configs.GetWeatherAPIKey() + "&q=" + strings.ReplaceAll(weather.City, " ", "_") + "&aqi=no")
	if error != nil {
		return error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return error
	}
	var w WeatherAPI
	error = json.Unmarshal(body, &w)
	if error != nil {
		return error
	}

	if (WeatherAPI{}) == w {
		return fmt.Errorf("can not find zipcode")
	}

	weather.FromCelsius(w.Current.TempC)

	return nil
}
