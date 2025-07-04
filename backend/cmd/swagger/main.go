package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Serve OpenAPI spec
	r.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../api/openapi.yaml")
	})

	// Serve Swagger UI
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../../api/swagger-ui/"))))

	log.Println("Swagger UI server starting on :8081")
	log.Println("Visit http://localhost:8081 to view the API documentation")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}