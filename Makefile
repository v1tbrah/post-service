up:
	docker-compose up -d

down:
	docker-compose down

protogen:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	cd ppbapi && \
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative post-service.proto

mockgen:
	go install github.com/vektra/mockery/v2@v2.24.0
	go generate ./...

migrateversiongen:
	cd migrations && \
	go generate ./...

test:
	go test ./...

test_with_db:
	docker-compose -f test.docker-compose.yaml up --build --abort-on-container-exit post-service-test
	docker-compose -f test.docker-compose.yaml down --volumes

push_to_dockerhub:
	docker build -t v1tbrah/post-service:v1 .
	docker tag v1tbrah/post-service:v1 v1tbrah/post-service:v1-release
	docker push v1tbrah/post-service:v1-release
