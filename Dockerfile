
FROM golang:1.20.5-alpine
WORKDIR /app
COPY . .
COPY go.mod go.sum ./
RUN go mod download
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
