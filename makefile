BINARY_NAME=main.out

build:
	go build -o ./out/bin/${BINARY_NAME} main.go

run:
	go build -o ./out/bin/${BINARY_NAME} main.go
	./out/bin/${BINARY_NAME} start

dep: 
	go mod download

clean:
	go mod tidy
	go clean

lint:
	golangci-lint run --enable-all
	