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
const requestNotError = "Request should not fail"

func setupWebServer() *WebServer {
	server := NewWebServer(testPort)
	server.setupDependencies()

	return server
}

func addHandlers(webServer *WebServer, routes []struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}) {
	for _, route := range routes {
		webServer.AddHandler(route.Path, route.Handler, route.Method)
	}
}

func startServer(t *testing.T, webServer *WebServer) func() {
	webServer.Start()

	go webServer.Run()
	time.Sleep(500 * time.Millisecond)

	return func() {
		err := webServer.Stop()
		assert.NoError(t, err)
	}
}

func performRequest(t *testing.T, method, endpoint string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, testBaseUrl+endpoint, nil)
	require.NoError(t, err, "Failed to create request")

	res, err := client.Do(req)

	return res, err
}

func TestAddHandler(t *testing.T) {
	webServer := setupWebServer()

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

func TestWebServerLifecycle(t *testing.T) {
	webServer := setupWebServer()

	routes := []struct {
		Path    string
		Method  string
		Handler http.HandlerFunc
	}{
		{
			Path:   testEndpoint,
			Method: "HEAD",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			Path:   testEndpoint,
			Method: "GET",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			Path:   testEndpoint,
			Method: "POST",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			},
		},
		{
			Path:   testEndpoint,
			Method: "PUT",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			Path:   testEndpoint,
			Method: "PATCH",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			Path:   testEndpoint,
			Method: "DELETE",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
		},
	}

	addHandlers(webServer, routes)
	defer startServer(t, webServer)()

	t.Run("Uptime must not be nil", func(t *testing.T) {
		assert.NotNil(t, webServer.UptimeService, "uptimeService should not be nil")
	})

	t.Run("Valid HEAD Handler", func(t *testing.T) {
		res, err := performRequest(t, "HEAD", testEndpoint)
		require.NoError(t, err, requestNotError)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Valid GET Handler", func(t *testing.T) {
		res, err := performRequest(t, "GET", testEndpoint)
		require.NoError(t, err, requestNotError)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Valid POST Handler", func(t *testing.T) {
		res, err := performRequest(t, "POST", testEndpoint)
		require.NoError(t, err, requestNotError)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run("Valid PUT Handler", func(t *testing.T) {
		res, err := performRequest(t, "PUT", testEndpoint)
		require.NoError(t, err, requestNotError)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Valid PATCH Handler", func(t *testing.T) {
		res, err := performRequest(t, "PATCH", testEndpoint)
		require.NoError(t, err, requestNotError)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Valid DELETE Handler", func(t *testing.T) {
		res, err := performRequest(t, "DELETE", testEndpoint)
		require.NoError(t, err, requestNotError)
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}

func TestAddHandlerWithDuplicateMethods(t *testing.T) {
	webServer := setupWebServer()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	webServer.AddHandler(testEndpoint, handler, "GET")
	webServer.AddHandler(testEndpoint, handler, "GET")

	key := testEndpoint + "_GET"
	assert.Contains(t, webServer.Handlers, key)
	assert.Len(t, webServer.Handlers, 1)
}

func TestInvalidMethods(t *testing.T) {
	webServer := setupWebServer()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	webServer.AddHandler(testEndpoint, handler, "GET")

	defer startServer(t, webServer)()

	t.Run("Invalid Method", func(t *testing.T) {
		res, err := performRequest(t, "INVALID", testEndpoint)
		require.NoError(t, err, requestNotError)
		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	})
}

func TestWebServerStop(t *testing.T) {
	t.Run("Stop After Start", func(t *testing.T) {
		webServer := setupWebServer()

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}
		webServer.AddHandler(testEndpoint, handler, "GET")

		stopServer := startServer(t, webServer)

		res, err := performRequest(t, "GET", testEndpoint)
		require.NoError(t, err, "Request should not fail before stopping the server")
		assert.Equal(t, http.StatusOK, res.StatusCode)

		stopServer()

		_, err = performRequest(t, "GET", testEndpoint)
		assert.Error(t, err, "Request should fail after server is stopped")
		assert.Contains(t, err.Error(), "refused", "Error should indicate connection refused")
	})

	t.Run("Stop Without Start", func(t *testing.T) {
		webServer := setupWebServer()
		webServer.Server = nil

		err := webServer.Stop()

		assert.NoError(t, err, "Calling Stop on a non-started server should not return an error")
	})
}

func TestWebServerErrorScenarios(t *testing.T) {
	t.Run("Start With Invalid Port", func(t *testing.T) {
		webServer := NewWebServer("invalidPort")
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}
		webServer.AddHandler(testEndpoint, handler, "GET")
		assert.Panics(t, func() {
			webServer.Start()
			webServer.Run()
		})
	})

	t.Run("Run Without Start", func(t *testing.T) {
		webServer := setupWebServer()
		assert.PanicsWithValue(t, "server not started: call Start() before Run()", func() {
			webServer.Run()
		})
	})

	t.Run("Run Start Without Error", func(t *testing.T) {
		webServer := setupWebServer()
		assert.NotPanics(t, func() {
			webServer.Start()
		})
	})
}
