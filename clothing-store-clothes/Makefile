generate:
		protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				pkg/pb/clothes.proto

run:
		go run ./cmd/clothes