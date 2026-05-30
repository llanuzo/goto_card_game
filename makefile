.PHONY: check run test docker genmocks

check: test genmocks
	go mod tidy
	staticcheck ./...
	gofmt -s -w .
	govulncheck ./...

run:
	go run -race ./cmd/card-game

	
test: 
	go test -timeout 30s -tags unit,store,integration ./...

docker:
	docker stop card-game || true
	docker rm card-game || true
	docker build -t card-game .
	docker run -d --name card-game -p 8080:8080 -p 10001:10001 card-game

genmocks:
	mockery --config .mockery.yml
	
