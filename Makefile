run:
	air

pactmode:
	PACT_MODE=true go run main.go

test:
	go test ./...

watch:
	gow test ./...

testv:
	go test -v ./...

docker:
	docker build -t book-lab-api:latest .

dockerpush:
	docker tag book-lab-api:latest gregbarozzi/book-lab-api:latest
	docker push gregbarozzi/book-lab-api:latest

dockerrun:
	docker run -it -p 8080:8080 --rm book-lab-api:latest

.PHONY: run pactmode docker dockerpush dockerrun test testv watch
