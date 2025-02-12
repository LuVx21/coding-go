ARG CR=${CR:-""}
ARG REGISTRY=${CR:+$CR/}
FROM ${REGISTRY}golang:1.24-alpine3.21 AS builder

RUN apk add --no-cache upx
WORKDIR /build

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux
COPY . .
RUN go mod download && go mod tidy && go build -ldflags="-w -s" -v -o luvx main.go && upx -9 luvx

# -------------------------------------------

FROM ${REGISTRY}alpine:3.21

WORKDIR /app
ENV log_LogDir=/root/data/logs TZ=Asia/Shanghai
COPY --from=builder /build/luvx .
COPY config/config-dev.yml config/config-dev.yml

# 将时区设置为东八区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo "Asia/Shanghai" > /etc/timezone # && apk del tzdata

EXPOSE 8090

CMD ["nohup", "./luvx", "&"]
