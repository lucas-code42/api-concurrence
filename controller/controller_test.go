package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucas-code42/api-race/models"
)

func TestGetMetrics(t *testing.T) {
	req, err := http.NewRequest(http.MethodTrace, "http://127.0.0.0.1:8080/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	getMetrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected code 200, but got %v", w.Code)
	}

	var expectedResponse models.Metrics
	if err := json.Unmarshal(w.Body.Bytes(), &expectedResponse); err != nil {
		t.Errorf("Expected: %v, Got: %v", expectedResponse, w.Body.String())
	}
}

func TestRouter(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.0.1:8080/64067030", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Router(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected code 200, but got %v", w.Code)
	}

	var expectedResponse models.ResponseDto
	if err := json.Unmarshal(w.Body.Bytes(), &expectedResponse); err != nil {
		t.Errorf("Expected: %v, Got: %v", expectedResponse, w.Body.String())
	}
}

func TestCepValidation(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.0.1:8080/6406703", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Router(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected code 500, but got %v", w.Code)
	}
}

func TestEnableCors(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, "http://127.0.0.0.1:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Router(w, req)
	t.Log(w.Body.String())
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected code 405, but got %v", w.Code)
	}

	var expectedResponse models.ResponseDto
	if err := json.Unmarshal(w.Body.Bytes(), &expectedResponse); err != nil {
		t.Errorf("Expected: %v, Got: %v", expectedResponse, w.Body.String())
	}
}

func TestGetMetricsRouterFlow(t *testing.T) {
	req, err := http.NewRequest(http.MethodTrace, "http://127.0.0.0.1:8080/64067030", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Router(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected code 200, but got %v", w.Code)
	}

	var expectedResponse models.ResponseDto
	if err := json.Unmarshal(w.Body.Bytes(), &expectedResponse); err != nil {
		t.Errorf("Expected: %v, Got: %v", expectedResponse, w.Body.String())
	}
}
