# ---- Build stage ----
FROM --platform=linux/amd64 golang:1.27 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- Run stage ----
FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
