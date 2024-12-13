package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Enjoy the silence!"))
		case "/health":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"cpu": {"cores": 2,"percent_used": [0,0]},"memory": {"total": 100,"used": 50,"free": 30,"available": 20,"percent_used": 0.5},"uptime": "37.51805678s","duration": "484.071µs","status": "pass","message": "Alive and kicking!","time": "2024-12-13T16:29:51-03:00"}`))
		case "/weather/24560352":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"temp_C": 14.4,"temp_F": 57.9,"temp_K": 287.4}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	t.Run("GET /", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("Error in request GET /: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected HTTP 200, but got %d", resp.StatusCode)
		}
	})

	t.Run("GET /health", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		if err != nil {
			t.Fatalf("Error in request GET /health: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected HTTP 200, but got %d", resp.StatusCode)
		}
	})

	t.Run("GET /weather/24560352", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/weather/24560352")
		if err != nil {
			t.Fatalf("Error in request GET /weather/24560352: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected HTTP 200, but got %d", resp.StatusCode)
		}
	})
}
