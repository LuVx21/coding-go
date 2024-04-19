# ARG IMAGE=golang:1.22
ARG IMAGE=golang:1.22-alpine3.19
FROM ${IMAGE} as pkg

WORKDIR /pkg
COPY go.mod ./
RUN go mod download && go mod verify

# -------------------------------------------

FROM pkg as builder

WORKDIR /build

COPY . .
RUN go mod download && go mod verify
# RUN go build -v -o . ./... # 所有main
RUN go build -v -o luvx main.go

# -------------------------------------------

FROM alpine:3.19
#FROM debian:bookworm

WORKDIR /app
COPY --from=builder /build/luvx .
COPY config/config-dev.yml config/config-dev.yml

EXPOSE 8090

CMD ["nohup", "./luvx", "&"]
