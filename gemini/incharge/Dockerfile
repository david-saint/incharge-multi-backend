FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /bin/incharge ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S incharge && adduser -S incharge -G incharge

WORKDIR /app

COPY --from=builder /bin/incharge .

USER incharge

EXPOSE 8080

ENTRYPOINT ["./incharge"]
