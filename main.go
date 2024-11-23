package main

import (
	"fmt"
	"net/http"
	"weatherzip/configs"
	"weatherzip/internal/infra/web/webserver"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	handlerRoot := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Enjoy the silence!"))
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webserver.AddHandler("/test", handler, "GET")
	webserver.AddHandler("/", handlerRoot, "GET")

	fmt.Println("Starting web server on port", cfg.WebServerPort)
	webserver.Start()
	webserver.Run()
}
