include .env
export

dev:
	go run main.go

migrate-down:
	migrate -database ${DATABASE_URL} -path migrations down

migrate:
	migrate -path migrations -database ${DATABASE_URL} up

dockerize:
	docker build -t codingpop/simplerest:v5 -f Dockerfile .
	docker push codingpop/simplerest:v5