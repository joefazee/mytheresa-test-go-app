package models

import (
	"fmt"
	"time"
)

type Product struct {
	ID       int
	Sku      string
	Name     string
	Price    int
	Discount int
	Category string
	Created  time.Time
}
type Products []*Product

type PublicProduct struct {
	Sku      string `json:"sku,omitempty"`
	Name     string `json:"name,omitempty"`
	Price    Price  `json:"price"`
	Category string `json:"category,omitempty"`
}

type Price struct {
	Original           int     `json:"original,omitempty"`
	Final              int     `json:"final,omitempty"`
	DiscountPercentage *string `json:"discount_percentage"`
	Currency           string  `json:"currency,omitempty"`
}

// Marshall encode to json
func (p *Product) Marshall() interface{} {

	var price = Price{
		Original: p.Price,
		Final:    p.Price - (p.Price / 100 * p.Discount),
		Currency: "EUR", // This should be moved to db
	}

	if p.Discount > 0 {
		var discountPercent = fmt.Sprintf("%d%%", p.Discount)
		price.DiscountPercentage = &discountPercent
	}

	return PublicProduct{
		Sku:      p.Sku,
		Name:     p.Name,
		Category: p.Category,
		Price:    price,
	}

}

func (products Products) Marshall() interface{} {

	result := make([]interface{}, len(products))
	for i, p := range products {
		result[i] = p.Marshall()
	}
	return result

}
