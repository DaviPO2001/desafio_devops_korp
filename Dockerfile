#Build da aplicação
FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

RUN go build -o http-server-projeto-korp .

#Execução da aplicação
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/http-server-projeto-korp .

EXPOSE 8080

CMD ["./http-server-projeto-korp"]
