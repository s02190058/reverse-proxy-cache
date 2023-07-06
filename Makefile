include .env
export

.PHONY: gen
gen:
	protoc --go_out=. --go_opt=paths=import \
	--go-grpc_out=. --go-grpc_opt=paths=import \
	api/thumbnail/v1/thumbnail.proto

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run:
	go run ./cmd/app/main.go

.PHONY: test
test:
	go test -v -race ./...

.PHONY: compose-up
compose-up:
	docker-compose up --build

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans