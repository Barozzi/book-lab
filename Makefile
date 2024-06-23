run:
	go run main.go

pactmode:
	PACT_MODE=true go run main.go

test:
	go test -v ./...

docker:
	docker build -t book-lab-api:latest .

dockerrun:
	docker run -it -p 8080:8080 --rm book-lab-api:latest

.PHONY: run pactmode docker dockerrun
