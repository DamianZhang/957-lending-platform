MIGRATE_PATH=repository/db/migration
DB_SOURCE=postgresql://root:secret@localhost:5432/957-lending-platform?sslmode=disable

postgres:
	docker run --name 957-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=957-lending-platform -d postgres:16-alpine

access_postgres:
	docker exec -i -t 957-postgres psql -U root 957-lending-platform

new_migration:
	migrate create -ext sql -dir "$(MIGRATE_PATH)" -seq "$(name)"

migrate_up:
	migrate -path "$(MIGRATE_PATH)" -database "$(DB_SOURCE)" -verbose up

migrate_up_1:
	migrate -path "$(MIGRATE_PATH)" -database "$(DB_SOURCE)" -verbose up 1

migrate_down:
	migrate -path "$(MIGRATE_PATH)" -database "$(DB_SOURCE)" -verbose down

migrate_down_1:
	migrate -path "$(MIGRATE_PATH)" -database "$(DB_SOURCE)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres access_postgres new_migration migrate_up migrate_up_1 migrate_down migrate_down_1 sqlc test
