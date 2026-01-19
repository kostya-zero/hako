default: run

run:
    go run . run

build:
    go build -o ./hako

build-exe:
    go build -o .\hako.exe
