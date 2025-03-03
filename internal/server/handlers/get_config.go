package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/lvestera/slot-machine/internal/models"
	"github.com/lvestera/slot-machine/internal/server/config"
)

type GetConfigHandler struct {
	Cfg *config.Config
}

func (h GetConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	coefficients := h.ComputeCoeff()
	resp, err := json.Marshal(coefficients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *GetConfigHandler) ComputeCoeff() []models.Coefficient {

	resultCoefficients := make([]models.Coefficient, 0, len(h.Cfg.Coefficients))
	var coefficient models.Coefficient

	for letter, values := range h.Cfg.Coefficients {
		distribution := math.Round(values[0] / values[1] * 10000)

		coefficient = models.Coefficient{
			Symbol:       letter,
			Distribution: distribution,
			Cost:         values[1],
		}

		resultCoefficients = append(resultCoefficients, coefficient)
	}

	return resultCoefficients
}
