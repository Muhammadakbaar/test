package search

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/akbar/k24/internal/product"
	"github.com/olivere/elastic/v7"
)

type SearchHandler struct {
	esClient  *elastic.Client
	indexName string
}

func NewSearchHandler(esClient *elastic.Client, indexName string) *SearchHandler {
	return &SearchHandler{esClient: esClient, indexName: indexName}
}
func (h *SearchHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	searchQuery := elastic.NewMultiMatchQuery(query, "product_name", "drug_generic", "company").
		Fuzziness("AUTO").
		Operator("and")

	searchResult, err := h.esClient.Search().
		Index(h.indexName).
		Query(searchQuery).
		Do(ctx)
	if err != nil {
		http.Error(w, "Error executing search", http.StatusInternalServerError)
		return
	}

	var products []product.Product
	for _, hit := range searchResult.Hits.Hits {
		var p product.Product
		if err := json.Unmarshal(hit.Source, &p); err != nil {
			// Handle error
			continue
		}
		p.ID = hit.Id
		p.Score = *hit.Score
		products = append(products, p)
	}

	sort.SliceStable(products, func(i, j int) bool {
		if products[i].Score != products[j].Score {
			return products[i].Score > products[j].Score
		}
		return strings.ToLower(products[i].ProductName) < strings.ToLower(products[j].ProductName)
	})

	response := map[string][]product.Product{"results": products}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
