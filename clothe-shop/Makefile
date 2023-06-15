auth-main:
		protoc --go_out=pkg/pb/auth --go_opt=paths=source_relative \
				--go-grpc_out=pkg/pb/auth --go-grpc_opt=paths=source_relative \
				pkg/pb/auth.proto

clothe-main:
		protoc --go_out=./clothe --go_opt=paths=source_relative \
				--go-grpc_out=./clothe --go-grpc_opt=paths=source_relative \
				pkg/pb/clothes.

