# ARG IMAGE=golang:1.23
ARG IMAGE=golang:1.23-alpine3.19
FROM ${IMAGE} as pkg

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# RUN apk add --no-cache gcc
WORKDIR /pkg
COPY go.mod ./
RUN go mod download && go mod verify

# -------------------------------------------

FROM pkg as builder

WORKDIR /build

COPY . .
RUN go mod download && go mod tidy
# RUN go build -v -o . ./... # 所有main
RUN go build -ldflags="-w -s" -v -o luvx main.go # && upx -9 luvx

# -------------------------------------------

FROM alpine:3.19
#FROM debian:bookworm

WORKDIR /app
ENV log_LogDir=/root/data/logs TZ=Asia/Shanghai
COPY --from=builder /build/luvx .
COPY config/config-dev.yml config/config-dev.yml

# 将时区设置为东八区
RUN echo "https://mirrors.aliyun.com/alpine/v3.19/main/" > /etc/apk/repositories \
    && echo "https://mirrors.aliyun.com/alpine/v3.19/community/" >> /etc/apk/repositories \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

EXPOSE 8090

CMD ["nohup", "./luvx", "&"]
