package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/example/newyear-days/internal/daysuntil"
)

const DateLayout = "2006-01-02"

type daysResponse struct {
	Date             string `json:"date"`
	DaysUntilNewYear int    `json:"days_until_new_year"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func DaysUntilNewYearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "only GET is supported")
		return
	}

	target, err := targetDate(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := daysResponse{
		Date:             target.Format(DateLayout),
		DaysUntilNewYear: daysuntil.DaysUntilNewYear(target),
	}
	writeJSON(w, http.StatusOK, resp)
}

func targetDate(r *http.Request) (time.Time, error) {
	raw := r.URL.Query().Get("date")
	if raw == "" {
		return time.Now().UTC(), nil
	}

	parsed, err := time.Parse(DateLayout, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid \"date\" parameter %q: expected format %s", raw, DateLayout)
	}
	return parsed, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("api: failed to encode JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
