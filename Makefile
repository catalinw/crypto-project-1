lint:
	golangci-lint run

build:
	go build ./...

test:
	go test ./...

coverage:
	go test -v -coverpkg=./... -coverprofile=coverage.cov ./... && go tool cover -html=coverage.cov

run:
	docker compose up --build

api-test:
	cd tests && godog

generate:
	go generate ./...