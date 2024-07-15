# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /rover cmd/rover/main.go

# runtime image
FROM alpine:3.18

WORKDIR /root/

COPY --from=build /rover .

CMD ["./rover"]
