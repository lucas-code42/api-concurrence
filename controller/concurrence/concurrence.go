package concurrence

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lucas-code42/api-race/models"
)

func getViaCep(cep string, viaCepChannel chan<- models.ViaCep) {
	url := strings.Replace(models.ViaCepUrl, "xxx", cep, 1)
	res, err := http.Get(url)
	if err != nil {
		log.Println("VIA CEP - error to mount request")
		// viaCepChannel <- models.ViaCep{Error: errors.New("error to mount request")}
		atomic.AddInt32(&models.ViaCepError, 1)
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println("VIA CEP - error to close body")
			// viaCepChannel <- models.ViaCep{Error: errors.New("error to close body")}
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("VIA CEP - error to read body")
		// viaCepChannel <- models.ViaCep{Error: errors.New("error to read body")}
		atomic.AddInt32(&models.ViaCepError, 1)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Println("VIA CEP - error statuscode")
		// viaCepChannel <- models.ViaCep{Error: errors.New("error statuscode")}
		atomic.AddInt32(&models.ViaCepError, 1)
		return
	}

	var responseModel models.ViaCep
	if err := json.Unmarshal(body, &responseModel); err != nil {
		log.Println("VIA CEP - error to unmarshal")
		// viaCepChannel <- models.ViaCep{Error: errors.New("error to unmarshal")}
		atomic.AddInt32(&models.ViaCepError, 1)
		return
	}
	viaCepChannel <- responseModel
}

func getBrasilAberto(cep string, brasilAbertoChannel chan<- models.BrasilAberto) {
	url := models.BrasilAbertoUrl + cep
	res, err := http.Get(url)
	if err != nil {
		log.Println("BRASILABERTO - error to mount request")
		// brasilAbertoChannel <- models.BrasilAberto{Error: errors.New("error to mount request")}
		atomic.AddInt32(&models.BrasilAbertoError, 1)
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println("BRASILABERTO - error to close body")
			// brasilAbertoChannel <- models.BrasilAberto{Error: errors.New("error to close body")}
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("BRASILABERTO - error to read body")
		// brasilAbertoChannel <- models.BrasilAberto{Error: errors.New("error to read body")}
		atomic.AddInt32(&models.BrasilAbertoError, 1)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Println("BRASILABERTO - error statuscode")
		// brasilAbertoChannel <- models.BrasilAberto{Error: errors.New("error statuscode")}
		atomic.AddInt32(&models.BrasilAbertoError, 1)
		return
	}

	var responseModel models.BrasilAberto
	if err := json.Unmarshal(body, &responseModel); err != nil {
		log.Println("BRASILABERTO - error to unmarshal")
		// brasilAbertoChannel <- models.BrasilAberto{Error: errors.New("error to unmarshal")}
		atomic.AddInt32(&models.BrasilAbertoError, 1)
		return
	}
	brasilAbertoChannel <- responseModel
}

func CepRace(cep string) models.ResponseDto {
	viaCepChannel := make(chan models.ViaCep)
	brasilAbertoChannel := make(chan models.BrasilAberto)

	go getBrasilAberto(cep, brasilAbertoChannel)
	go getViaCep(cep, viaCepChannel)

	select {
	case viaCepResponse := <-viaCepChannel:
		atomic.AddInt32(&models.ViaCepTotal, 1)
		log.Println("** VIA CEP WINS **\n", viaCepResponse)
			cepDto := models.CepDto{
				Cep:         viaCepResponse.Cep,
				Logradouro:  viaCepResponse.Logradouro,
				Complemento: viaCepResponse.Complemento,
				Bairro:      viaCepResponse.Bairro,
				Localidade:  viaCepResponse.Localidade,
				Uf:          viaCepResponse.Uf,
				ApiOrigin:   "viaCep",
			}
			return models.ResponseDto{Data: cepDto, Error: models.Err{ErrorMessage: nil}}

	case brasilAbertoResponse := <-brasilAbertoChannel:
		atomic.AddInt32(&models.BrasilAbertoTotal, 1)
		log.Println("** BRASIL ABERTO WINS **\n", brasilAbertoResponse)
		if brasilAbertoResponse.Error == nil {
			cepDto := models.CepDto{
				Cep:         brasilAbertoResponse.Result.Street,
				Logradouro:  brasilAbertoResponse.Result.District,
				Complemento: brasilAbertoResponse.Result.Complement,
				Bairro:      brasilAbertoResponse.Result.StateShortname,
				Localidade:  brasilAbertoResponse.Result.State,
				Uf:          brasilAbertoResponse.Result.City,
				ApiOrigin:   "brasilAberto",
			}
			return models.ResponseDto{Data: cepDto, Error: models.Err{ErrorMessage: nil}}
		}
		atomic.AddInt32(&models.BrasilAbertoError, 1)
		log.Println("BRASIL ABERTO ERROR ->", brasilAbertoResponse.Error)

	case <-time.After(time.Second * 3):
		atomic.AddInt32(&models.TimeOut, 1)
		log.Println("** TIME OUT **")
	}

	return models.ResponseDto{}
}
