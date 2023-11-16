package concurrence

import (
	"testing"

	"github.com/lucas-code42/api-race/models"
	"github.com/stretchr/testify/assert"
)

func TestGetViaCep(t *testing.T) {
	viaCepChannel := make(chan models.ViaCep, 1)

	t.Run("ViaCep", func(t *testing.T) {
		getViaCep("64067030", viaCepChannel)
		result := <-viaCepChannel
		assert.Equal(t, "64067-030", result.Cep)
	})
}

func TestGetBrasilAberto(t *testing.T) {
	// remember to use chan size
	brasilAbertoChannel := make(chan models.BrasilAberto, 1)

	// Teste de sucesso
	t.Run("BrasilAberto", func(t *testing.T) {
		getBrasilAberto("64067030", brasilAbertoChannel)
		result := <-brasilAbertoChannel
		assert.Equal(t, "64067030", result.Result.Zipcode)
	})

}

func TestCepRace(t *testing.T) {
	const (
		ViaCepOrigin       = "viaCep"
		BrasilAbertoOrigin = "brasilAberto"
	)
	t.Run("CepRace", func(t *testing.T) {
		chResponseDto := CepRace("64067030")

		if chResponseDto.Data.ApiOrigin != ViaCepOrigin && chResponseDto.Data.ApiOrigin != BrasilAbertoOrigin {
			t.Log(chResponseDto.Data.ApiOrigin)
			t.Error("cepRace error")
		}

		t.Log(chResponseDto.Data.ApiOrigin, "wins!")
	})
}
