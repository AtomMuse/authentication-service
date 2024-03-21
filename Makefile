run-auth:
	go run cmd/authentication/main.go

test-coverage:
	mkdir -p coverage
	go test -race -short -v -coverprofile coverage/cover.out ./...
	go tool cover -html=coverage/cover.out

gen-swag:
	swag init -d ./cmd/authentication,./handler/authentication -o ./cmd/authentication/doc --pd