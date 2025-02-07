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

	r1 := h.Cfg.WinAAAFreq / h.Cfg.WinAAACost
	r2 := h.Cfg.WinBBBFreq / h.Cfg.WinBBBCost
	r3 := h.Cfg.WinCCCFreq / h.Cfg.WinCCCCost

	d1 := math.Round(math.Cbrt(r1)*1000) / 1000
	d2 := math.Round(math.Cbrt(r2)*1000) / 1000
	d3 := math.Round(math.Cbrt(r3)*1000) / 1000

	coeff := make([]models.Coefficient, 0, 3)

	coeffA := models.Coefficient{
		Symbol:       "A",
		Distribution: d1,
		Cost:         h.Cfg.WinAAACost,
	}
	coeffB := models.Coefficient{
		Symbol:       "B",
		Distribution: d2,
		Cost:         h.Cfg.WinBBBCost,
	}
	coeffC := models.Coefficient{
		Symbol:       "C",
		Distribution: d3,
		Cost:         h.Cfg.WinCCCCost,
	}

	return append(coeff, coeffA, coeffB, coeffC)
}
