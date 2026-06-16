FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/crrc-pm-backend .

FROM alpine:3.20

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/crrc-pm-backend /app/crrc-pm-backend

EXPOSE 8080

USER app

CMD ["/app/crrc-pm-backend"]
