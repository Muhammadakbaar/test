package main

import (
	"context"
	"log"

	"github.com/akbar/k24/internal/elasticsearch"
	"github.com/akbar/k24/pkg/utils"
	"github.com/olivere/elastic/v7"
)

const (
	indexName = "products"
	excelPath = "Skill test_data_product.xlsx"
)

func main() {
	esClient, err := elasticsearch.NewClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	if err := elasticsearch.CreateIndex(esClient, indexName); err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	products, err := utils.ReadProductsFromExcel(excelPath)
	if err != nil {
		log.Fatalf("Error reading products from Excel: %s", err)
	}

	bulkRequest := esClient.Bulk()
	for _, p := range products {
		req := elastic.NewBulkIndexRequest().Index(indexName).Id(p.ID).Doc(p)
		bulkRequest = bulkRequest.Add(req)
	}

	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		log.Fatalf("Error executing bulk request: %s", err)
	}

	if bulkResponse.Errors {
		log.Println("Bulk request had errors")
	}

	log.Printf("Indexed %d products", len(bulkResponse.Items))
}
