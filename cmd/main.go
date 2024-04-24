package main

import (
	"fmt"

	"github.com/tiagocosta/cloud-run/configs"
	"github.com/tiagocosta/cloud-run/internal/infra/web"
	"github.com/tiagocosta/cloud-run/internal/infra/web/webserver"
)

func main() {
	webserver := webserver.NewWebServer(configs.GetWebServerPort())
	webWeatherHandler := web.NewWebWeatherHandler()
	webserver.AddHandler("/weather", webWeatherHandler.Get)
	fmt.Println("Starting web server on port", configs.GetWebServerPort())
	webserver.Start()
}
