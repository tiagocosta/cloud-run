package main

import (
	"fmt"

	"github.com/tiagocosta/cloud-run/configs"
	"github.com/tiagocosta/cloud-run/internal/infra/web"
	"github.com/tiagocosta/cloud-run/internal/infra/web/webserver"
)

func main() {
	cfg := configs.LoadConfig[configs.Conf](".")

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webWeatherHandler := web.NewWebWeatherHandler()
	webserver.AddHandler("/weather", webWeatherHandler.Get)
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	webserver.Start()
}
