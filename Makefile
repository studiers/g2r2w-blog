build: build-server build-client

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/blog.proto

build-server: gen-proto
	go build ./server

build-client: gen-proto
	go build ./client
