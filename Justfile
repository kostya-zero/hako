set shell := ["bash", "-c"]
set windows-shell := ["pwsh.exe", "-NoLogo", "-Command"]

binaryPath := if os() == "windows" { './build/hako.exe' } else { './build/hako' }

# Runs build recipe
default: build

# Update dependencies
update:
    go get -u ./...

# Build the project to an executable
build:
    go build -o {{ binaryPath }} ./cmd/hako/

# Run the application with optional arguments
run *ARGS:
    go run ./cmd/hako/ {{ ARGS }}
