package server

import (
	"fmt"
	"net/http"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/geocode", GeocodeHandler)
	mux.HandleFunc("/forecast", ForecastHandler)
	mux.HandleFunc("/map/tiles/", TileHandler)

	return mux
}

func Start() {
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", NewMux())
}
