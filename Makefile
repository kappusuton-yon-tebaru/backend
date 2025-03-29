.PHONY: gen dev-agent dev-backend dev-consumer lint build

gen:
	@wire ./cmd/agent/agent/
	@wire ./cmd/backend/backend/
	@wire ./cmd/consumer/consumer/
	@wire ./cmd/sidecarlogger/sidecarlogger/

docs:
	@swag fmt
	@swag init --parseDependency -g ./cmd/backend/main.go -o ./cmd/backend/docs

neogen:
	@go run tools/wire/main.go \
		./cmd/agent/agent \
		./cmd/backend/backend \
		./cmd/consumer/consumer

dev-agent:
	@air -c ./cmd/agent/.air.toml

dev-backend:
	@air -c ./cmd/backend/.air.toml

dev-consumer:
	@air -c ./cmd/consumer/.air.toml

build:
	@docker build -t public.ecr.aws/r2n4f6g5/agent:latest -f cmd/agent/Dockerfile .
	@docker build -t public.ecr.aws/r2n4f6g5/consumer:latest -f cmd/consumer/Dockerfile .
	@docker build -t public.ecr.aws/r2n4f6g5/sidecarlogger:latest -f cmd/sidecarlogger/Dockerfile .

push:
	@docker push public.ecr.aws/r2n4f6g5/agent:latest
	@docker push public.ecr.aws/r2n4f6g5/consumer:latest
	@docker push public.ecr.aws/r2n4f6g5/sidecarlogger:latest

lint:
	@golangci-lint run --timeout=5m

apply:
	@kubectl apply --server-side --field-manager=system -f deployment/master.local.yaml

delete:
	@kubectl delete -f deployment/master.local.yaml

manifest:
	@go run ./tools/manifest/main.go
