package main

type config struct {
	port int
	env  string // Environment (development|staging|production)
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}
