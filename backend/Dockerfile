FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/server/main.go


# Build the final image
FROM alpine:3.20.1 AS app

WORKDIR /app
COPY --from=build /app/main /app/main
CMD ["./main"]
