install:
	@ dep ensure

run:
	@ docker-compose up -d
	@ go run main.go

test:
	@ go test -race ./...