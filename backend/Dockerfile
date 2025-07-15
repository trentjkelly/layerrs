# Building backend (Go) binary
FROM golang:1.23-alpine
WORKDIR /dist
COPY . .
RUN go build -o api ./cmd/api
EXPOSE 8080
CMD ["./api"]
