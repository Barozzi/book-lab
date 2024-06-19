run:
	go run main.go

docker:
	docker build -t book-lab-api:latest .

dockerrun:
	docker run -it -p 8080:8080 --rm book-lab-api:latest

.PHONY: run docker
