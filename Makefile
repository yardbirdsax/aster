coverage:
	go tool cover -func=coverage.out
generate:
	go generate ./...
test:
	go test -v -count=1 -coverprofile coverage.out -race ./...
tidy:
	go mod tidy