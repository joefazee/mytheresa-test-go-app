package mock

import (
	"github.com/joefazee/mytheresa/pkg/data"
	"github.com/joefazee/mytheresa/pkg/models"
)

var products models.Products = models.Products{
	&models.Product{
		ID:       1,
		Sku:      "000001",
		Name:     "BV Lean leather ankle boots",
		Category: "boots",
		Price:    89000,
		Discount: 30,
	},

	&models.Product{
		ID:       2,
		Sku:      "000002",
		Name:     "BV Lean leather ankle boots",
		Category: "boots",
		Price:    99000,
		Discount: 0,
	},

	&models.Product{
		ID:       3,
		Sku:      "000003",
		Name:     "Ashlington leather ankle boots",
		Category: "boots",
		Price:    71000,
		Discount: 0,
	},

	&models.Product{
		ID:       4,
		Sku:      "000004",
		Name:     "Naima embellished suede sandals",
		Category: "sandals",
		Price:    79500,
		Discount: 0,
	},

	&models.Product{
		ID:       5,
		Sku:      "000005",
		Name:     "Nathane leather sneakers",
		Category: "sneakers",
		Price:    59000,
		Discount: 0,
	},
}

type ProductModel struct {
}

func (m ProductModel) GetAll(filter data.Filter) (models.Products, data.Metadata, error) {

	var result models.Products
	switch {

	case filter.Category != "" && filter.PriceLessThan > 0:
		for _, p := range products {
			if p.Category == filter.Category && p.Price <= filter.PriceLessThan {
				result = append(result, p)
			}
		}
		return result, data.Metadata{}, nil
	case filter.Category != "" && filter.PriceLessThan <= 0:
		for _, p := range products {
			if p.Category == filter.Category {
				result = append(result, p)
			}
		}
		return result, data.Metadata{}, nil
	case filter.PriceLessThan > 0 && filter.Category == "":
		for _, p := range products {
			if p.Price <= filter.PriceLessThan {
				result = append(result, p)
			}
		}
		return result, data.Metadata{}, nil
	}

	return products, data.Metadata{}, nil
}
