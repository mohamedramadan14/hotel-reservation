build: 
	@go build -o bin/api

seed:
	@go run scripts/seed.go

run: build
	@./bin/api

lint:
	@ golangci/golangci-lint-action@v3
test:
	@go test -v ./... --count=1