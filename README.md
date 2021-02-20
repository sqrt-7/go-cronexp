# CRON Expression Parser

## Description
This application takes a single cron expression string as a command-line argument and prints out
the expanded view of the given schedule.

## Requirements
- `go 1.14+`
- `make`

## Setup
1. Clone git repository to local machine
```bash
git clone https://github.com/sqrt-7/go-cronexp.git && cd go-cronexp
```
2. Download dependencies
```bash
go mod download
```
3. Compile application
```bash
make build
```
4. Run with cron expression as the argument
```bash
./run_cronexp "*/15 0 1,15 * 1-5 /usr/bin/find"
```

## Testing
- Unit tests are located at `pkg/cronexp/cronexp_test.go`
- Run unit tests:
```bash
make unit-test
```

## Troubleshooting
- If the application has to be compiled for a specific os/architecture:
```bash
# For valid GOOS/GOARCH combinations run 
go tool dist list

# Replace GOOS and GOARCH with the chosen values
GOOS=darwin GOARCH=amd64 go build -o ./run_cronexp ./cmd/cli/main.go
```