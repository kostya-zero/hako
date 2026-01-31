$CGO_ENABLED=0
$GOEXPERIMENT=greenteagc

default: run

run:
    go run cmd/main/main.go run

build:
    go build -o hako ./cmd/main

build-exe:
    go build -o hako.exe ./cmd/main
