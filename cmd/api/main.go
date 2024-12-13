package main

import (
	"fmt"
	"net/http"

	"github.com/vs0uz4/weatherzip/configs"
	"github.com/vs0uz4/weatherzip/internal/infra/web"
	"github.com/vs0uz4/weatherzip/internal/infra/web/webserver"
	"github.com/vs0uz4/weatherzip/internal/service"
	"github.com/vs0uz4/weatherzip/internal/usecase"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{}

	cpuService := service.NewCPUService()
	memoryService := service.NewMemoryService()
	uptimeService := service.NewUptimeService()
	cepService := service.NewCepService(httpClient, cfg.CepAPIUrl)
	weatherService := service.NewWeatherService(httpClient, cfg.WeatherAPIUrl, cfg.WeatherAPIKey, cfg.WeatherAPILanguage)

	healthCheckUseCase := usecase.NewHealthCheckUseCase(cpuService, memoryService, uptimeService)
	wheaterByCepUseCase := usecase.NewWeatherByCepUsecase(cepService, weatherService)

	handlerRoot := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Enjoy the silence!")); err != nil {
			http.Error(w, "Unable to write response", http.StatusInternalServerError)
		}
	}

	handlerHealth := web.NewHealthHandler(healthCheckUseCase).GetHealth
	handlerWeather := web.NewWeatherHandler(wheaterByCepUseCase).GetWeatherByCep

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webserver.AddHandler("/weather/{cep}", handlerWeather, "GET")
	webserver.AddHandler("/health", handlerHealth, "GET")
	webserver.AddHandler("/", handlerRoot, "GET")

	fmt.Println("Starting web server on port", cfg.WebServerPort)
	webserver.Start()
	webserver.Run()
}
