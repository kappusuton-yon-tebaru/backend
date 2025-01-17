.PHONY: gen

gen:
	@echo "Generating agent..."
	@wire ./cmd/agent/agent

	@echo "Generating backend..."
	@wire ./cmd/backend/backend
