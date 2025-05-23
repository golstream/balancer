FROM golang:1.22-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -o /app/bin/app cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=build /app/bin/app .

CMD ["./app"]
