package server

import (
	"fmt"
	"net/http"
)

func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request for: " + r.URL.Path)
		fmt.Fprintf(w, "Endpoint reached: %s", r.URL.Path)
	})

	http.HandleFunc("/geocode", GeocodeHandler)

	http.HandleFunc("/forecast", ForecastHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
