### Update proto
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative post-service.proto
```

### Update mocks
```
mockery --name PostServiceClient
```
