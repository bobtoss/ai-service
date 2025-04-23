FROM golang:1.20-alpine as builder

WORKDIR /build
COPY . /build
RUN go build -o app .

FROM alpine:3.18.0 as hoster

WORKDIR /app

COPY --from=builder /build/app ./app
COPY --from=builder /build/config ./config 

ENTRYPOINT ["./app"]
