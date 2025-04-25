include .env
export $(shell sed 's/=.*//' .env)

DB_PG="postgres://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

# SQLi

sqli:
	xo --src ../../xo --out ./db --schema public schema ${DB_PG}

# CRM DB

migrate-up:
	cd ./migrations && goose postgres ${DB_PG} up

migrate-down:
	cd ./migrations && goose postgres ${DB_PG} down

migrate-down-to-zero:
	cd ./migrations && goose postgres ${DB_PG} down-to 0

migrate-reup:
	make migrate-down-to-zero
	make migrate-up

migrate-create:
	cd ./migrations && goose create $(name) sql

# Tests

prepare:
	docker-compose -f docker-compose.yaml up -d
	docker logs database 2>&1 | grep -q "database system is ready to accept connections"
	make migrate-reup
	make sqli

test:
	make prepare
	grc go test . -count=1