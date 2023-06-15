FROM golang:1.18-alpine

WORKDIR /app/clothing-store-brands
COPY . .

RUN go mod download

CMD ["go", "run", "cmd/brands/brands.go"]