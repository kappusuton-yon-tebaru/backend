.PHONY: gen dev-agent dev-backend dev-builder-consumer lint build

gen:
	@wire ./cmd/agent/agent/
	@wire ./cmd/backend/backend/
	@wire ./cmd/builder-consumer/builderconsumer/

neogen:
	@go run tools/wire/main.go \
		./cmd/agent/agent \
		./cmd/backend/backend \
		./cmd/builder-consumer/builderconsumer

dev-agent:
	@air -c ./cmd/agent/.air.toml

dev-backend:
	@air -c ./cmd/backend/.air.toml

dev-builder-consumer:
	@air -c ./cmd/builder-consumer/.air.toml

lint:
	@golangci-lint run

apply:
	@kubectl apply --server-side --field-manager=system -f deployment/master.yaml

delete:
	@kubectl delete -f deployment/master.yaml
