package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joefazee/mytheresa/migration"
	"github.com/joefazee/mytheresa/pkg/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "secret"
	dbName   = "celeritas_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var testDB *sql.DB
var resource *dockertest.Resource
var pool *dockertest.Pool
var productsModel ProductModel

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_TYPE", "postgres")
	os.Setenv("UPPER_DB_LOG", "ERROR")

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		err = pool.Purge(resource)
		log.Fatal("error cleaning resources: ", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to docker: %s", err)
	}

	err = migration.Run(testDB)
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	productsModel = ProductModel{DB: testDB}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestProductModel_GetAll(t *testing.T) {

	testTable := []struct {
		name      string
		filter    data.Filter
		wantCount int
		wantError bool
	}{
		{name: "no category and priceLessThan", filter: data.Filter{Category: "", PriceLessThan: 0, PageSize: 5, Page: 1}, wantCount: 5, wantError: false},
		{name: "filter by category", filter: data.Filter{Category: "boots", PageSize: 5, Page: 1}, wantCount: 3, wantError: false},
		{name: "filter by priceLessThan", filter: data.Filter{PriceLessThan: 80000, PageSize: 5, Page: 1}, wantCount: 3, wantError: false},
		{name: "filter by priceLessThan", filter: data.Filter{Category: "boots", PriceLessThan: 80000, PageSize: 5, Page: 1}, wantCount: 1, wantError: false},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			products, _, err := productsModel.GetAll(tt.filter)
			gotCount := len(products)

			if err != nil {
				t.Errorf("error running query; %q", err)
			}

			if tt.wantError && err == nil {
				t.Errorf("expected error; got nil")
			}

			if gotCount != tt.wantCount {
				t.Errorf("expected %d of products; got %d", tt.wantCount, gotCount)
			}

		})
	}

}
