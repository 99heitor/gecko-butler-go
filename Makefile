export GO111MODULE=on
BINARY_NAME=bot

all: deps build
install:
	go install bot.go
build:
	go build -o $(BINARY_NAME) main.go
static-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY_NAME} .
test:
	go test -v ./...
clean:
	go clean
	rm -f $(BINARY_NAME)
deps:
	go build -v ./...
upgrade:
	go get -u