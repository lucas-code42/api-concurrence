package main

import (
	"net/http"

	"github.com/lucas-code42/api-race/controller"
)

func main() {
	http.HandleFunc("/", controller.Router)
	http.ListenAndServe(":8080", nil)
}
