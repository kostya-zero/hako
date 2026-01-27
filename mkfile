$CGO_ENABLED=0
$GOEXPERIMENT=greenteagc

default: run

run:
    go run . run

build:
    go build -o ./hako

build-exe:
    go build -o .\hako.exe
