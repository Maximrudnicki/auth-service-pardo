
build:
	protoc -Iproto --go_opt=module=auth_service --go_out=. --go-grpc_opt=module=auth_service --go-grpc_out=. proto/*.proto
	go build -o bin/auth_service.exe ./cmd/.