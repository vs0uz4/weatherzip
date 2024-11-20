package webserver

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewWebServer(t *testing.T) {
	port := ":8080"
	webServer := NewWebServer(port)

	assert.Equal(t, port, webServer.WebServerPort)
	assert.NotNil(t, webServer.Router)
	assert.NotNil(t, webServer.Handlers)
}

func TestAddHandler(t *testing.T) {
	port := ":8080"
	webServer := NewWebServer(port)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler("/test", handler, "GET")

	key := "/test_GET"

	assert.Contains(t, webServer.Handlers, key)
	assert.Len(t, webServer.Handlers, 1)
	reflect.DeepEqual(handler, webServer.Handlers[key].Handler)
	assert.Equal(t, "GET", webServer.Handlers[key].Method)
}

func TestStart(t *testing.T) {
	port := ":8080"
	webServer := NewWebServer(port)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler("/test", handler, "GET")

	go webServer.Start()

	time.Sleep(100 * time.Millisecond)

	defer func() {
		webServer.Router = chi.NewRouter()
	}()

	req, _ := http.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	webServer.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
