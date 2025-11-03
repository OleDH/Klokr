FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o klokr .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/klokr .
EXPOSE 8080
CMD ["./klokr", "-serve"]

