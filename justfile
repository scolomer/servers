gen:
    buf generate

[working-directory: 'go']
server:
    go run grpc/server/main.go

[working-directory: 'go']
client:
    go run grpc/client/main.go
