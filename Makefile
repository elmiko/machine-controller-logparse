all: build

build:
	go build -o machine-controller-logparse .

test:
	go test -v ./...
