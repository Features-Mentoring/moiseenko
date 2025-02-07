package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lvestera/slot-machine/internal/models"
	"github.com/lvestera/slot-machine/internal/storage"
)

type GetResultHandler struct {
	Db storage.DBRepository
}

func (h GetResultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var results map[int]models.Result
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &results)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = h.Db.AddBatch(results)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
