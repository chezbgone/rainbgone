package server

import (
	"fmt"
	"net/http"
)

func Start() {
	http.HandleFunc("/geocode", GeocodeHandler)

	http.HandleFunc("/forecast", ForecastHandler)

	http.HandleFunc("/map/tiles/", TileHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
