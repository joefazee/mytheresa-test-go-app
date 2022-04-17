package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joefazee/mytheresa/migration"
	"github.com/joefazee/mytheresa/pkg/logger"
	"github.com/joefazee/mytheresa/pkg/models/postgres"
	_ "github.com/lib/pq"
	"os"
)

func main() {

	var cfg config
	cfg.db.dsn = os.Getenv("DATABASE_DSN")

	flag.IntVar(&cfg.port, "port", 5555, "application listening port")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-cons", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-cons", 25, "PostgreSQL max idle connections")

	flag.Parse()

	// create logger
	appLogger := logger.New(os.Stdout, logger.LevelInfo)

	// SET UP DB CONNECTION
	db, err := connectToDatabase(cfg.db.dsn)
	if err != nil {
		appLogger.Fatal(err, map[string]string{"mysql": "unable to connect to mysql", "dsn": cfg.db.dsn})
	}
	defer db.Close()
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// run migration
	err = migration.Run(db)
	if err != nil {
		appLogger.Error(err, map[string]string{"mysql": "unable to execute migration"})
	}

	// init app
	app := &application{
		config:   cfg,
		products: &postgres.ProductModel{DB: db},
		logger:   appLogger,
	}

	app.logger.Info("starting application server", map[string]string{"port": fmt.Sprintf("%d", app.config.port)})
	if err := app.serve(); err != nil {
		app.logger.Fatal(err, nil)
	}

}

func connectToDatabase(dsn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
