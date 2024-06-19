# Dockerfile.distroless
FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./book-lab-api .

EXPOSE 8080

CMD ["./book-lab-api"]
