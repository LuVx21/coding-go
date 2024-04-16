FROM golang:1.22-alpine3.19 as builder
#FROM golang:1.22 as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# RUN go build -v -o . ./... # 所有main
RUN go build -v -o luvx main.go

FROM alpine:3.19
#FROM debian:bookworm

WORKDIR /app
COPY --from=builder /build/luvx .
COPY config/config-dev.yml config/config-dev.yml

EXPOSE 8090

CMD ["nohup", "./luvx", "&"]
