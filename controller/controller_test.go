package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucas-code42/api-race/models"
	"github.com/stretchr/testify/mock"
)

type MockFuncs interface {
	enableCors(w http.ResponseWriter, r *http.Request)
	getCep(w http.ResponseWriter, r *http.Request)
	getMetrics(w http.ResponseWriter, r *http.Request)
	response(w http.ResponseWriter, r *http.Request)
}

type TestMock struct {
	mock.Mock
}

func (m *TestMock) enableCors(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *TestMock) getCep(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *TestMock) getMetrics(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *TestMock) response(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestGetMetrics(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.0.1:8080/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	getMetrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Esperava um código de status 200, mas recebeu %v", w.Code)
	}

	expectedBody := models.Metrics{}
	if err := json.Unmarshal(w.Body.Bytes(), &expectedBody); err != nil {
		t.Errorf("O corpo da resposta não corresponde ao esperado. Esperado: %v, Recebido: %v", expectedBody, w.Body.String())
	}
}

func TestRouter(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.0.1:8080/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Router(w, req)

	t.Log(w.Body.String())
}
