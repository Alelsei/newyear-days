package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDaysUntilNewYearHandler_Integration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(DaysUntilNewYearHandler))
	defer server.Close()

	t.Run("explicit date returns the correct count and 200 OK", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/?date=2024-12-31")
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body struct {
			Date             string `json:"date"`
			DaysUntilNewYear int    `json:"days_until_new_year"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if body.Date != "2024-12-31" {
			t.Errorf("date = %q, want %q", body.Date, "2024-12-31")
		}
		if body.DaysUntilNewYear != 1 {
			t.Errorf("days_until_new_year = %d, want %d", body.DaysUntilNewYear, 1)
		}
	})

	t.Run("leap year date returns the correct count", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/?date=2024-01-01")
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body struct {
			DaysUntilNewYear int `json:"days_until_new_year"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if body.DaysUntilNewYear != 366 {
			t.Errorf("days_until_new_year = %d, want %d", body.DaysUntilNewYear, 366)
		}
	})

	t.Run("missing date defaults to today and still succeeds", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		var body struct {
			Date             string `json:"date"`
			DaysUntilNewYear int    `json:"days_until_new_year"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if body.Date != time.Now().UTC().Format(DateLayout) {
			t.Errorf("date = %q, want today's date %q", body.Date, time.Now().UTC().Format(DateLayout))
		}
		if body.DaysUntilNewYear <= 0 || body.DaysUntilNewYear > 366 {
			t.Errorf("days_until_new_year = %d, want a value in (0, 366]", body.DaysUntilNewYear)
		}
	})

	t.Run("malformed date is rejected with 400 Bad Request", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/?date=not-a-date")
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}

		var body struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if body.Error == "" {
			t.Error("expected a non-empty error message")
		}
	})

	t.Run("unsupported method is rejected with 405 Method Not Allowed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, server.URL+"/", nil)
		if err != nil {
			t.Fatalf("failed to build request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
		}
	})
}
