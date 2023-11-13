package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/lucas-code42/api-race/controller/concurrence"
	"github.com/lucas-code42/api-race/models"
)

func getMetrics(w http.ResponseWriter, r *http.Request) {
	byteResponse, err := json.Marshal(models.Metrics{
		Total:   models.Total,
		TimeOut: models.TimeOut,
		ViaCep: struct {
			Total int32
			Ok    int32
			Error int32
		}{Total: models.ViaCepTotal, Ok: models.ViaCepOk, Error: models.ViaCepError},
		BrasilAberto: struct {
			Total int32
			Ok    int32
			Error int32
		}{Total: models.BrasilAbertoTotal, Ok: models.BrasilAbertoOk, Error: models.BrasilAbertoError},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to return metrics"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteResponse)
}

func response(w *http.ResponseWriter, r *http.Request, data []byte, statusCode int) {
	if len(data) == 0 {
		(*w).WriteHeader(statusCode)
		(*w).Write(data)
		return
	}
	(*w).WriteHeader(statusCode)
	(*w).Write(data)
}

func getCep(w *http.ResponseWriter, r *http.Request) {
	rawCep := strings.ReplaceAll(r.URL.Path, "/", "")
	proceed := func() bool {
		if _, err := strconv.Atoi(rawCep); err != nil || len(rawCep) < 8 {
			return false
		}
		return true
	}()
	if !proceed {
		response(w, r, []byte("error to mount ResponseDTO fora"), http.StatusInternalServerError)
		return
	}

	res := concurrence.CepRace(rawCep)
	atomic.AddInt32(&models.Total, 1)

	if res.Data.ApiOrigin == "brasilAberto" {
		atomic.AddInt32(&models.BrasilAbertoOk, 1)
	} else if res.Data.ApiOrigin == "viaCep" {
		atomic.AddInt32(&models.ViaCepOk, 1)
	}

	if res == (models.ResponseDto{}) {
		response(w, r, []byte("error to mount ResponseDTO fora"), http.StatusInternalServerError)
		return
	}
	byteResponse, err := json.Marshal(&res)
	if err != nil {
		response(w, r, []byte("error to mount ResponseDTO fora"), http.StatusInternalServerError)
		return
	}
	response(w, r, byteResponse, http.StatusInternalServerError)
}

func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func Router(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
	if r.Method == http.MethodGet {
		getCep(&w, r)
		return
	} else if r.Method == http.MethodTrace {
		getMetrics(w, r)
		return
	}

	dtoRes := models.ResponseDto{
		Data:  models.CepDto{},
		Error: models.Err{ErrorMessage: errors.New("method not allow, try a GET")},
	}
	res, err := json.Marshal(&dtoRes)
	if err != nil {
		response(&w, r, []byte("error to mount ResponseDTO"), http.StatusInternalServerError)
	}

	response(&w, r, res, http.StatusMethodNotAllowed)
}
