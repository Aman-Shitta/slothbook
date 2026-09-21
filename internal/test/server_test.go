package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aman-Shitta/slotbook/internal/server"
	"github.com/stretchr/testify/assert"
)

type HealthzResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}
}

func TestHealthzRoute(t *testing.T) {

	router := server.SetupRouter("8000")

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(
		"GET",
		"/healthz",
		nil,
	)

	router.ServeHTTP(w, req)

	// check status is 200
	assert.Equal(t, 200, w.Result().StatusCode)

	var response HealthzResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)

	// check status is healthy
	assert.Equal(t, "healthy", response.Data.Status)

	// check version is v0.0.1
	assert.Equal(t, "v0.1.0", response.Data.Version)

}
