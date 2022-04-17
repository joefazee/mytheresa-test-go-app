package postgres

import (
	"database/sql"
	"github.com/joefazee/mytheresa/pkg/data"
	"github.com/joefazee/mytheresa/pkg/models"
)

const (
	querySelectProductsNoQuery                    = "SELECT COUNT(*) OVER() total_records, id, sku, name, price, discount, category, created FROM products LIMIT $1 OFFSET $2"
	querySelectProductsByCategory                 = "SELECT COUNT(*) OVER() total_records, id, sku, name, price, discount, category, created FROM products WHERE category = $1 LIMIT $2 OFFSET $3"
	querySelectProductsByCategoryAndPriceLessThan = "SELECT COUNT(*) OVER() total_records, id, sku, name, price, discount, category, created FROM products WHERE category = $1 AND price <= $2 LIMIT $3 OFFSET $4"
	querySelectProductsByPressLessThan            = "SELECT COUNT(*)  OVER() total_records, id, sku, name, price, discount, category, created FROM products WHERE price <= $1 LIMIT $2 OFFSET $3"
	queryInsertProduct                            = "INSERT INTO products (sku, name, price, discount, category, created) VALUES ($1, $2, $3, $4, $5, UTC_TIMESTAMP())"
)

// ProductModel  allows us to configure dependencies
type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) GetAll(filter data.Filter) (models.Products, data.Metadata, error) {

	var err error
	var rows *sql.Rows

	switch {

	case filter.Category != "" && filter.PriceLessThan > 0: // I was thinking of setting default value to -1
		rows, err = m.DB.Query(querySelectProductsByCategoryAndPriceLessThan, filter.Category, filter.PriceLessThan, filter.Limit(), filter.Offset())
	case filter.Category != "" && filter.PriceLessThan <= 0:
		rows, err = m.DB.Query(querySelectProductsByCategory, filter.Category, filter.Limit(), filter.Offset())
	case filter.PriceLessThan > 0 && filter.Category == "":
		rows, err = m.DB.Query(querySelectProductsByPressLessThan, filter.PriceLessThan, filter.Limit(), filter.Offset())
	default:
		rows, err = m.DB.Query(querySelectProductsNoQuery, filter.Limit(), filter.Offset()) // I could have used QueryContext instead
	}

	if err != nil {
		return nil, data.Metadata{}, err
	}
	defer rows.Close()

	var products []*models.Product
	totalRecords := 0

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.Sku,
			&product.Name,
			&product.Price,
			&product.Discount,
			&product.Category,
			&product.Created)

		if err != nil {
			return nil, data.Metadata{}, err
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)
	return products, metadata, nil

}
