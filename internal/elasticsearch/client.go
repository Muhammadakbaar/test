package elasticsearch

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
)

func NewClient() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		return nil, err
	}

	log.Println("Elasticsearch client initialized")
	return client, nil
}

func CreateIndex(client *elastic.Client, indexName string) error {
	ctx := context.Background()
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		mapping := `
		{
			"mappings": {
				"properties": {
					"product_name": {
						"type": "text",
						"analyzer": "standard"
					},
					"drug_generic": {
						"type": "text",
						"analyzer": "standard"
					},
					"company": {
						"type": "text",
						"analyzer": "standard"
					}
				}
			}
		}`
		_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
		log.Printf("Index %s created", indexName)
	} else {
		log.Printf("Index %s already exists", indexName)
	}

	return nil
}
