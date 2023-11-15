package models

// ResponseDto is the default Response from service
type ResponseDto struct {
	Data  CepDto `json:"data"`
	Error Err    `json:"error"`
}

// Err is the default error msg
type Err struct {
	ErrorMessage string `json:"errorMessage"`
}

// CepDto is the default DTO from service
type CepDto struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	ApiOrigin   string `json:"apiOrigin"`
}

// BrasilAberto is the struct that https://brasilaberto.com/ uses
type BrasilAberto struct {
	Meta struct {
		CurrentPage  int `json:"currentPage"`
		ItemsPerPage int `json:"itemsPerPage"`
		TotalOfItems int `json:"totalOfItems"`
		TotalOfPages int `json:"totalOfPages"`
	} `json:"meta"`
	Result struct {
		Street         string `json:"street"`
		Complement     string `json:"complement"`
		District       string `json:"district"`
		DistrictID     int    `json:"districtId"`
		City           string `json:"city"`
		CityID         int    `json:"cityId"`
		IbgeID         int    `json:"ibgeId"`
		State          string `json:"state"`
		StateShortname string `json:"stateShortname"`
		Zipcode        string `json:"zipcode"`
	} `json:"result"`
	Error error
}

// ViaCep is the struct that https://viacep.com.br/ uses
type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Error       error
}

// InitProviders return's the initalize providers url
var (
	ViaCepUrl       = "https://viacep.com.br/ws/xxx/json/"
	BrasilAbertoUrl = "https://brasilaberto.com/api/v1/zipcode/"
)

// Metrics
var (
	Total int32

	TimeOut int32

	ViaCepTotal int32
	ViaCepOk    int32
	ViaCepError int32

	BrasilAbertoTotal int32
	BrasilAbertoOk    int32
	BrasilAbertoError int32
)

type Metrics struct {
	Total int32

	TimeOut int32

	ViaCep struct {
		Total int32
		Ok    int32
		Error int32
	}

	BrasilAberto struct {
		Total int32
		Ok    int32
		Error int32
	}
}
