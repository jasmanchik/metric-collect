proto:
	protoc -I=./api/grpc --go_out=./api/grpc/go --go_opt=paths=source_relative --go-grpc_out=./api/grpc/go --go-grpc_opt=paths=source_relative ./api/grpc/metrics.proto

run: lint start

lint:
	golangci-lint run

start:
	go run ./cmd/metrics -config=./config/local_test.yaml