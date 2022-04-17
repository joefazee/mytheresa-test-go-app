.PHONY: run/api
run/api:
	@echo 'Make sure you have set DATABASE_DSN environmental variable using postgres format'
	@go run ./cmd/api


.PHONY: run/test
run/test:
	@echo 'running unit and integration test (Note: Integration test uses docker. )'
	@go test -v ./...



.PHONY: docker/run
docker/run: docker/stop
	@echo 'running via docker '
	@docker compose up

.PHONY: docker/stop
docker/stop:
	@docker compose down



