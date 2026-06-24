package server

import (
	"fmt"
	"log"
	"net/http"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/geocode", GeocodeHandler)
	mux.HandleFunc("/forecast", ForecastHandler)

	return mux
}

func Start() {
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", NewMux()))
}
