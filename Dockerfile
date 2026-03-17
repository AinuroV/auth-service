
FROM golang:1.24-alpine AS builder


RUN apk add --no-cache git

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

COPY --from=builder /app/auth-service .

COPY --from=builder /app/.env.example .

RUN chown -R appuser:appgroup /home/appuser


USER appuser


EXPOSE 8080


CMD ["./auth-service"]