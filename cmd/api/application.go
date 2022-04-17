package main

import (
	"github.com/joefazee/mytheresa/pkg/logger"
)

type application struct {
	config   config
	products ProductsInterface
	logger   *logger.Logger
}
