.PHONY: gen dev-agent dev-backend lint

gen:
	@go run tools/wire/main.go \
		./cmd/agent/agent \
		./cmd/backend/backend \
		./cmd/builder-consumer/builderconsumer

dev-agent:
	@air -c ./cmd/agent/.air.toml

dev-backend:
	@air -c ./cmd/backend/.air.toml

lint:
	@golangci-lint run
