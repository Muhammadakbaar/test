package utils

import (
	"fmt"

	"github.com/akbar/k24/internal/product"
	"github.com/xuri/excelize/v2"
)

func ReadProductsFromExcel(filePath string) ([]product.Product, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var products []product.Product
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			return nil, fmt.Errorf("row %d has fewer than 4 columns", i+1)
		}
		products = append(products, product.Product{
			ID:          row[0],
			ProductName: row[1],
			DrugGeneric: row[2],
			Company:     row[3],
		})
	}

	return products, nil
}
