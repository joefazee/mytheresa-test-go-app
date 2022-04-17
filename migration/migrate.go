package migration

import (
	"database/sql"
	"log"
	"strings"
)

// For large tables, this scriot should be managed by a migration library
var sqlStr = `
DROP TABLE IF EXISTS products;

CREATE TABLE IF NOT EXISTS  products (
    id bigserial PRIMARY KEY,
    sku TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    price INT NOT NULL DEFAULT 0,
    discount SMALLINT NOT NULL DEFAULT 0 CHECK(discount BETWEEN 0 AND 100),
    category VARCHAR(255) NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

INSERT INTO products (id, sku, name, price, discount, category, created)
VALUES (1, '000001', 'BV Lean leather ankle boots', 89000, 30, 'boots', '2022-04-16 19:14:20'),
       (2, '000002', 'BV Lean leather ankle boots', 99000, 0, 'boots', '2022-04-16 19:14:20'),
       (3, '000003', 'Ashlington leather ankle boots', 71000, 0, 'boots', '2022-04-16 19:14:20'),
       (4, '000004', 'Naima embellished suede sandals', 79500, 10, 'sandals', '2022-04-16 19:14:20'),
       (5, '000005', 'Nathane leather sneakers', 59000, 0, 'sneakers', '2022-04-16 19:14:20');
`

func Run(db *sql.DB) error {

	queries := strings.Split(sqlStr, ";")
	if len(queries) > 0 {
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if len(query) > 0 {
				_, err := db.Exec(query)
				if err != nil {
					log.Fatal(err)
				}
			}

		}

	}

	return nil
}
