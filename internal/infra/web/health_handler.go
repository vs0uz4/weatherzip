package web

import (
	"encoding/json"
	"net/http"

	"github.com/vs0uz4/weatherzip/internal/usecase"
)

type HealthHandler struct {
	useCase usecase.HealthCheckUseCase
}

func NewHealthHandler(u usecase.HealthCheckUseCase) *HealthHandler {
	return &HealthHandler{u}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	health, err := h.useCase.GetHealth()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(health)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
