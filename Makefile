build:
	@go build -o bin/halo-suster cmd/server/main.go

run: build
	@./bin/halo-suster