# The Single Endpoint Product App in Go!

The entire application is based on docker. I have tested it with docker on the mac.

You can access the application on http://localhost:5555/v1/products

## How to run the application

```shell
git clone git@github.com:joefazee/mytheresa-test-go-app.git  && cd mytheresa-test-go-app && make docker/run
```

### Important commands

```shell
make run/test
# Runs both unit and integration test
```

```shell
make docker/run
# Runs the application on http://localhost:5555/v1/products
```

```shell
make docker/stop
# Stop the application
```

## Running outside docker

Running the application outside docker requires more setup. You will need to have go go1.18 and postgres installed.

set DATABASE_DSN envriomental variable

```shell
export DATABASE_DSN="host=localhost port=5432 user=postgres password=password dbname=dbname sslmode=disable timezone=UTC connect_timeout=5"
```

```shell
make run/api
```

### Some of my design decisions

- Structure the application to grow into a medium project if need be.

- Limit the usage of external libraries: I want to make sure I show more original codes. For large projects, we sure
  need to use libraries like Gin, Gorm, etc.

- Keep the code transparent. No magic.
- Keep the folder structure self-explanatory and straightforward.
- I used Go instead PHP just to write more code.
- Test critical parts. Areas like the handlers, db access are all tested.
