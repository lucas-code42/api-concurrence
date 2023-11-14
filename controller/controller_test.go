package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucas-code42/api-race/models"
	"github.com/stretchr/testify/mock"
)

type TestMock struct {
	mock.Mock
}

type MockFuncs interface {
	enableCors(w http.ResponseWriter, r *http.Request)
	getCep(w http.ResponseWriter, r *http.Request)
	getMetrics(w http.ResponseWriter, r *http.Request) error
	response(w http.ResponseWriter, r *http.Request)
}

func (m *TestMock) enableCors(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *TestMock) getCep(w *http.ResponseWriter, r *http.Request) {
	m.Called(&w, r)
}

func (m *TestMock) getMetrics(w http.ResponseWriter, r *http.Request) error {
	args := m.Called(&w, r)
	return args.Error(0)
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
		t.Errorf("Expected code 200, but got %v", w.Code)
	}

	expectedBody := models.Metrics{}
	if err := json.Unmarshal(w.Body.Bytes(), &expectedBody); err != nil {
		t.Errorf("Expected: %v, Got: %v", expectedBody, w.Body.String())
	}
}

func TestRouter(t *testing.T) {
	tm := &TestMock{}
	req, err := http.NewRequest("GET", "http://127.0.0.0.1:8080/64067030", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// Espera que getCep seja chamado exatamente uma vez
	tm.On("getCep", w, mock.Anything).Return(nil).Once()

	Router(w, req)

	// Verifica se todas as expectativas foram atendidas
	tm.AssertExpectations(t)
}
