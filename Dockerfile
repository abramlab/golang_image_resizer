FROM golang:1.18.5-alpine3.15 as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./image-resizer ./
RUN mkdir -p images

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/image-resizer /app/image-resizer
COPY --from=builder /build/images /app/images
ENTRYPOINT ["/app/image-resizer", "--input", "/app/images"]