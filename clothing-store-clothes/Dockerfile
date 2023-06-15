FROM golang:1.18-alpine

WORKDIR /app/clothing-store-clothes
COPY . .

RUN go mod download

CMD ["go", "run", "cmd/clothes/clothes.go"]