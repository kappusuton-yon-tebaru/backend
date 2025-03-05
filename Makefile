.PHONY: gen dev-agent dev-backend dev-builder-consumer lint build

gen:
	@wire ./cmd/agent/agent/
	@wire ./cmd/backend/backend/
	@wire ./cmd/builder-consumer/builderconsumer/
	@swag init -g ./cmd/backend/main.go -o cmd/backend/docs

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

build:
	@docker build -t public.ecr.aws/r2n4f6g5/agent -f cmd/agent/Dockerfile .
	@docker build -t public.ecr.aws/r2n4f6g5/builder-consumer -f cmd/builder-consumer/Dockerfile .

push:
	@docker push public.ecr.aws/r2n4f6g5/agent
	@docker push public.ecr.aws/r2n4f6g5/builder-consumer

lint:
	@golangci-lint run --timeout=5m

apply:
	@kubectl apply --server-side --field-manager=system -f deployment/master.local.yaml

delete:
	@kubectl delete -f deployment/master.local.yaml

manifest:
	@go run ./tools/manifest/main.go
