package webserver

import "net/http"

type WebServerInterface interface {
	AddHandler(path string, handler http.HandlerFunc, method string)
	Start()
	Run()
	Stop() error
}
