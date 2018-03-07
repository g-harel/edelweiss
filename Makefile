install:
	@ dep ensure

run:
	@ docker-compose up -d
	@ go run ./services/gateway/main.go

test:
	@ go test -race ./...