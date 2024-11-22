package webserver

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testBaseUrl = "http://localhost:8080"
const testEndpoint = "/test"
const testPort = ":8080"

func TestNewWebServer(t *testing.T) {
	webServer := NewWebServer(testPort)

	assert.Equal(t, testPort, webServer.WebServerPort)
	assert.NotNil(t, webServer.Router)
	assert.NotNil(t, webServer.Handlers)
}

func TestAddHandler(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")

	key := testEndpoint + "_GET"

	assert.Contains(t, webServer.Handlers, key)
	assert.Len(t, webServer.Handlers, 1)
	reflect.DeepEqual(handler, webServer.Handlers[key].Handler)
	assert.Equal(t, "GET", webServer.Handlers[key].Method)
}

func TestAddHandlerDuplicate(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")
	webServer.AddHandler(testEndpoint, handler, "GET")

	key := testEndpoint + "_GET"

	assert.Contains(t, webServer.Handlers, key)
	assert.Len(t, webServer.Handlers, 1)
	reflect.DeepEqual(handler, webServer.Handlers[key].Handler)
	assert.Equal(t, "GET", webServer.Handlers[key].Method)
}

func TestStart(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")

	webServer.Start()

	go webServer.Run()
	time.Sleep(500 * time.Millisecond)

	defer func() {
		err := webServer.Stop()
		assert.NoError(t, err)
	}()

	client := &http.Client{}

	req1, _ := http.NewRequest("GET", testBaseUrl+testEndpoint, nil)
	res1, err := client.Do(req1)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res1.StatusCode)
}

func TestAddHandlerWithDifferentMethods(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler1 := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	handler2 := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}

	webServer.AddHandler(testEndpoint, handler1, "GET")
	webServer.AddHandler(testEndpoint, handler2, "POST")

	key1 := testEndpoint + "_GET"
	key2 := testEndpoint + "_POST"

	assert.Contains(t, webServer.Handlers, key1)
	assert.Contains(t, webServer.Handlers, key2)
	assert.Len(t, webServer.Handlers, 2)
	reflect.DeepEqual(handler1, webServer.Handlers[key1].Handler)
	reflect.DeepEqual(handler2, webServer.Handlers[key2].Handler)
	assert.Equal(t, "GET", webServer.Handlers[key1].Method)
	assert.Equal(t, "POST", webServer.Handlers[key2].Method)
}

func TestAddHandlerWithInvalidMethod(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "INVALID")

	key := testEndpoint + "_INVALID"

	assert.Contains(t, webServer.Handlers, key)
	assert.Len(t, webServer.Handlers, 1)
	reflect.DeepEqual(handler, webServer.Handlers[key].Handler)
	assert.Equal(t, "INVALID", webServer.Handlers[key].Method)
}

func TestStartWithMultipleHandlers(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler1 := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	handler2 := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}
	handler3 := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	webServer.AddHandler(testEndpoint, handler1, "HEAD")
	webServer.AddHandler(testEndpoint, handler1, "GET")
	webServer.AddHandler(testEndpoint, handler2, "POST")
	webServer.AddHandler(testEndpoint, handler1, "PUT")
	webServer.AddHandler(testEndpoint, handler1, "PATCH")
	webServer.AddHandler(testEndpoint, handler3, "DELETE")

	webServer.Start()

	go webServer.Run()
	time.Sleep(500 * time.Millisecond)

	defer func() {
		err := webServer.Stop()
		assert.NoError(t, err)
	}()

	client := &http.Client{}

	req1, _ := http.NewRequest("GET", testBaseUrl+testEndpoint, nil)
	res1, err := client.Do(req1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res1.StatusCode)

	req2, _ := http.NewRequest("POST", testBaseUrl+testEndpoint, nil)
	res2, err := client.Do(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res2.StatusCode)

	req3, _ := http.NewRequest("PUT", testBaseUrl+testEndpoint, nil)
	res3, err := client.Do(req3)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res3.StatusCode)

	req4, _ := http.NewRequest("PATCH", testBaseUrl+testEndpoint, nil)
	res4, err := client.Do(req4)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res4.StatusCode)

	req5, _ := http.NewRequest("DELETE", testBaseUrl+testEndpoint, nil)
	res5, err := client.Do(req5)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res5.StatusCode)

	req6, _ := http.NewRequest("HEAD", testBaseUrl+testEndpoint, nil)
	res6, err := client.Do(req6)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res6.StatusCode)
}

func TestStartWithInvalidMethod(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")

	webServer.Start()

	go webServer.Run()
	time.Sleep(500 * time.Millisecond)

	defer func() {
		err := webServer.Stop()
		assert.NoError(t, err)
	}()

	client := &http.Client{}

	req1, _ := http.NewRequest("HEAD", testBaseUrl+testEndpoint, nil)
	res1, err := client.Do(req1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res1.StatusCode)
}

func TestStartWithInvalidPort(t *testing.T) {
	webServer := NewWebServer("invalidPort")

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")

	assert.Panics(t, func() {
		webServer.Start()
		webServer.Run()
	})
}

func TestStop(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")

	webServer.Start()
	go webServer.Run()

	time.Sleep(500 * time.Millisecond)

	err := webServer.Stop()
	assert.NoError(t, err)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", testBaseUrl+testEndpoint, nil)
	_, err = client.Do(req)

	assert.Error(t, err)
}

func TestStopWithoutStart(t *testing.T) {
	webServer := NewWebServer(testPort)

	err := webServer.Stop()
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	webServer := NewWebServer(testPort)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")

	webServer.Start()

	go func() {
		webServer.Run()
	}()

	time.Sleep(500 * time.Millisecond)

	client := &http.Client{}

	req, err := http.NewRequest("GET", testBaseUrl+testEndpoint, nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	err = webServer.Stop()
	assert.NoError(t, err)
}

func TestRunWithoutStart(t *testing.T) {
	webServer := NewWebServer(testPort)

	assert.Panics(t, func() {
		webServer.Run()
	})
}
