FROM golang:1.23-alpine AS builder

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o bin/application cmd/server/main.go

FROM alpine:3.21 AS runner

WORKDIR /opt

COPY --from=builder /opt/bin/application ./

CMD ["./application"]
