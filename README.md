# Crypto-API


## Run the application

In order to run the application, docker needs to be installed on your machine.
Run `make run` to run crypto-api and the postgres database in docker.

## Run tests

There are 2 types of tests written for the API: API tests that call the endpoints and integration tests for the challenge service.

Api tests:
- cucumber tests written using [godog](https://github.com/cucumber/godog)
- these tests will call the real endpoints in the crypto-API
- in order to run the tests you have to run the application first using `make run`
- run API tests using:
    - `go install github.com/cucumber/godog/cmd/godog@v0.12.0` - just to be sure godog is installed  
    - `make api-test`

Challenge service tests:
- integration tests for the challenge service 
- challenge repository is mocked using [mockgen](https://github.com/golang/mock)
- business logic and expected calls to repository dependency are validated  
- run integration tests `make test`

## Mocks

Dependency mocks generated by [mockgen](https://github.com/golang/mock)

If the dependencies change, re-generate the mocks by running command `make generate`

## Send requests to the application

Please import postman collection `welthee.postman_collection.json` in order to call crypto-API endppints.

## How to use crypto-cli to generate signed tokens

Crypto-cli application can be used to create a token that contain a nonce using ES256 signature algorithm

In order to crypto-cli it please run:
`cd crypto-cli`
`./crypto-cli jwt <YOUR_NONCE_GENERATED_BY_CRYPTO_API>`

In order to build the crypto-cli application please run:
`cd crypto-cli`
`make build`


