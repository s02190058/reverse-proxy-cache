.PHONY: gen
gen:
	protoc --go_out=. --go_opt=paths=import \
	--go-grpc_out=. --go-grpc_opt=paths=import \
	api/thumbnail/v1/thumbnail.proto

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -v -race ./...