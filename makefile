include .env
export

dev:
	go run main.go

migrate:
	migrate -database ${DATABASE_URL} -path migrations down
	migrate -path migrations -database ${DATABASE_URL} up
	