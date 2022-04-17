package main

import (
	"github.com/joefazee/mytheresa/pkg/data"
	"github.com/joefazee/mytheresa/pkg/models"
)

type ProductsInterface interface {
	GetAll(filter data.Filter) (models.Products, data.Metadata, error)
}
