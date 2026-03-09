.PHONY: run wire swagger lint fmt build

run:
	go run main.go

wire:
	cd pkg/di && wire

swagger:
	swag init -g main.go -o docs/

lint:
	go vet ./...

fmt:
	go fmt ./...

build:
	go build -o cloudops main.go

docker-env:
	docker compose -f docker-compose-env.yaml up -d
