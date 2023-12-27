FROM golang:1-alpine AS builder

RUN mkdir /httpGoPing
WORKDIR /httpGoPing

COPY go.mod go.mod
COPY go.sum go.sum
COPY httpGoPing.go httpGoPing.go

RUN go build -o httpGoPing ./httpGoPing.go

FROM alpine:latest

COPY --from=builder /httpGoPing/httpGoPing /usr/bin/
EXPOSE 8080
USER nobody
CMD ["/usr/bin/httpGoPing"]
