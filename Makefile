.PHONY: gen dev-agent dev-backend

gen:
	@echo "Generating agent..."
	@wire ./cmd/agent/agent

	@echo "Generating backend..."
	@wire ./cmd/backend/backend

	@echo "Generating builder consumer..."
	@wire ./cmd/builder-consumer/builderconsumer/

dev-agent:
	@air -c ./cmd/agent/.air.toml

dev-backend:
	@air -c ./cmd/backend/.air.toml
