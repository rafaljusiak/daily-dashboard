FROM golang:1.22.2-bookworm

RUN go install github.com/air-verse/air@latest

RUN air init

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download && go mod verify

WORKDIR /app/src

CMD ["air", "-c", ".air.toml"]
