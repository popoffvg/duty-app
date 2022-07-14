build:
	docker-compose build duty-app

run:
	docker-compose up duty-app

test:
	go test -v ./...

migrate-up:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable' down

swag:
	swag init -g cmd/main.go