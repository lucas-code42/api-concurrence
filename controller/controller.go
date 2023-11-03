package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lucas-code42/api-race/models"
)

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
	}

	dtoRes := models.ResponseDto{
		Data:  models.CepDto{},
		Error: models.Err{ErrorMessage: "Method not allow, try a GET"},
	}
	res, err := json.Marshal(&dtoRes)
	if err != nil {
		response(&w, r, []byte("error to mount ResponseDTO"), http.StatusInternalServerError)
	}

	response(&w, r, res, http.StatusMethodNotAllowed)
}

func response(w *http.ResponseWriter, r *http.Request, data []byte, statusCode int) {
	if len(data) == 0 {
		(*w).WriteHeader(statusCode)
		(*w).Write(data)
		return
	}

}

func getCep(w *http.ResponseWriter, r *http.Request) {
	res := cepRace()
	fmt.Println(res)
}

func cepRace() models.ResponseDto {
	viaCepChannel := make(chan models.ViaCep)
	// brasilAbertoChannel := make(chan models.BrasilAberto)

	err := getViaCep("123", viaCepChannel)
	if err != nil {
		
	}

	select {
	case viaCepResponse := <-viaCepChannel:

	}

	return models.ResponseDto{}
}

func getViaCep(cep string, viaCepChannel chan<- models.ViaCep) error {
	url := strings.Replace(models.ViaCepUrl, "xxx", cep, 1)

	res, err := http.Get(url)
	if err != nil {
		return errors.New("error to mount request")
	}
	res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("error to read io")
	}

	if res.StatusCode != http.StatusOK {
		// return models.ViaCep{}, errors.New("error bad request")
		return errors.New("error bad request")
	}

	var responseModel models.ViaCep
	if err := json.Unmarshal(body, &responseModel); err != nil {
		return errors.New("error to unmarshal")
	}

	viaCepChannel <- responseModel
	return nil
}
