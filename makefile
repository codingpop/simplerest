include .env
export

dev:
	echo ${DATABASE_URL}
	go run main.go

migrate:
	migrate -database ${DATABASE_URL} -path migrations down
	migrate -path migrations -database ${DATABASE_URL} up

dockerize:
	docker build -t codingpop/simplerest:v3 -f Dockerfile .
	docker push codingpop/simplerest:v3