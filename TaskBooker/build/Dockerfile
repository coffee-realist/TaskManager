FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /booker ./cmd/main.go

FROM scratch

COPY --from=builder /booker /booker

EXPOSE 8000

ENTRYPOINT ["/booker"]
