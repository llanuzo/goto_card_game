.PHONY: check run docker genmocks

check: test
	go mod tidy
	staticcheck ./...
	gofmt -s -w .
	govulncheck ./...

run:
	go run -race ./cmd/card-game

	
test: 
	go test -timeout 30s -tags unit,store,integration ./...

docker:
	docker build -t card-game .
	docker run -p 8080:8080 card-game

# Will fail if there are no packages specified in .mockery.yml
genmocks:
	mockery --config .mockery.yml
	
