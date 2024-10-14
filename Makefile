generate:
		protoc -I protoss/proto protoss/proto/auth_service/auth_service.proto --go_out=protoss/gen/go --go_opt=paths=source_relative --go-grpc_out=protoss/gen/go/ --go-grpc_opt=paths=source_relative

run:
		go run auth_service/cmd/auth_serivce/main.go --config=auth_service/config/local.yaml

migrations:
		go run auth_service/cmd/migrator/main.go --storage-path=./auth_service/storage/auth_serivce.db --migrations-path=./auth_service/migrations
