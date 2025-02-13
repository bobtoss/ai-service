FROM golang:1.20-alpine as builder
WORKDIR /build
COPY . /build
RUN go build -o ai-service .

FROM alpine:3.18.0 as hoster
COPY --from=builder /build/.env* ./.env
COPY --from=builder /build/ai-service ./ai-service

# executable
ENTRYPOINT [ "./ai-service" ]