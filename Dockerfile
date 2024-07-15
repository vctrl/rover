ARG GOLANG_VERSION
ARG ALPINE_VERSION

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /rover cmd/rover/main.go

# runtime image
FROM alpine:${ALPINE_VERSION}

WORKDIR /root/

COPY --from=build /rover .

CMD ["./rover"]
