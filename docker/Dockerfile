FROM golang:1.24-alpine AS builder

WORKDIR /doc_server
COPY go.* ./
RUN go mod download

COPY ../.. ./

RUN go build -ldflags="-s -w" -v -o monolith ./cmd/doc_server

FROM alpine:latest AS runtime

COPY --from=builder /doc_server/monolith /doc_server/monolith

CMD ["/doc_server/monolith"]
