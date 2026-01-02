package main

import (
	"log"
	"net/http"

	"github.com/akbar/k24/internal/elasticsearch"
	"github.com/akbar/k24/internal/search"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/akbar/k24/docs"
)

const (
	indexName = "products"
)

func main() {
	esClient, err := elasticsearch.NewClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	if err := elasticsearch.CreateIndex(esClient, indexName); err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	r := mux.NewRouter()

	searchHandler := search.NewSearchHandler(esClient, indexName)
	r.HandleFunc("/search", searchHandler.SearchProducts).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
