package configs

import "github.com/spf13/viper"

var cfg *conf

type conf struct {
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
}

func init() {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
}

func GetWeatherAPIKey() string {
	return cfg.WeatherAPIKey
}

func GetWebServerPort() string {
	return cfg.WebServerPort
}
