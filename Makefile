lint:
	golangci-lint run

build:
	go build ./...

test:
	go test ./...

coverage:
	go test -v -coverpkg=./... -coverprofile=coverage.cov ./... && go tool cover -html=coverage.cov

run:
	docker compose up

api-tests:
	cd tests && godog