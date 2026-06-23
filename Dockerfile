FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o olive-blog-api

FROM alpine:latest
RUN apk update --no-cache && apk add --no-cache tzdata ca-certificates
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/olive-blog-api /app/olive-blog-api
EXPOSE 8080
CMD ["/app/olive-blog-api"]