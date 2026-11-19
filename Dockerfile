FROM golang:1.23-alpine AS builder
WORKDIR /build

COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/hello-go .

FROM alpine:3.24
RUN adduser -D -H appuser
WORKDIR /app
COPY --from=builder /out/hello-go ./hello-go

USER appuser
ENTRYPOINT ["./hello-go"]
