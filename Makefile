DB_URL=postgresql://postgres:postgres@localhost:5432/ppob?sslmode=disable

network:
	docker network create ppob-network

postgres:
	docker run --name postgres-local --network ppob-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb:
	docker exec -it postgres-local createdb --username=postgres --owner=postgres ppob

dropdb:
	docker exec -it postgres-local dropdb --username=postgres ppob

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ariandi/ppob_go/db/sqlc Store

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock