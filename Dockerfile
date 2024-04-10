FROM golang:1.22.2-bookworm

RUN go install github.com/cosmtrek/air@latest

RUN air init

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

CMD ["air", "-c", ".air.toml"]
