package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vs0uz4/weatherzip/internal/domain"
	"github.com/vs0uz4/weatherzip/internal/infra/web/webserver/middleware"
	"github.com/vs0uz4/weatherzip/internal/usecase/contracts"

	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	Usecase contracts.WeatherByCepUsecase
}

func NewWeatherHandler(uc contracts.WeatherByCepUsecase) *WeatherHandler {
	return &WeatherHandler{Usecase: uc}
}

func (h *WeatherHandler) GetWeatherByCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	weather, err := h.Usecase.GetWeatherByCep(cep)
	if err != nil {
		if errors.Is(err, domain.ErrZipcodeNotFound) {
			if rr, ok := w.(*middleware.ResponseRecorder); ok {
				rr.WriteError("Zipcode not found")
			}
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		if err.Error() == "invalid zipcode" {
			if rr, ok := w.(*middleware.ResponseRecorder); ok {
				rr.WriteError("Invalid zipcode")
			}
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if rr, ok := w.(*middleware.ResponseRecorder); ok {
			rr.WriteError("Internal server error")
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"temp_C": weather.Current.TempC,
		"temp_F": weather.Current.TempF,
		"temp_K": weather.Current.TempK,
	}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
