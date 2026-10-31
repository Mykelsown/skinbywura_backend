.PHONY: run test vet build seed import

run:
	go run ./cmd/api

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -o bin/api ./cmd/api

# Loads the 8 demo products. Use only on a fresh local database.
seed:
	psql "$$DATABASE_URL" -f db/seed.sql

# Loads or updates products from catalog.csv (copy catalog.example.csv first).
import:
	go run ./cmd/importproducts -file catalog.csv
