build: clean
	go build -o ./run_cronexp ./cmd/cli/main.go

clean:
	rm -f ./run_cronexp

unit-test:
	go test -v ./pkg/...