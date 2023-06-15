FROM golang:1.18-alpine

WORKDIR /app/clothe-shop
COPY . .

RUN go mod download

CMD ["go", "run", "cmd/clothing-store/main.go"]