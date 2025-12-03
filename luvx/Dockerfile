ARG CR=${CR:-""}
ARG REGISTRY=${CR:+$CR/}
FROM ${REGISTRY}golang:1.25-alpine3.22 AS builder

COPY --from=registry.cn-shanghai.aliyuncs.com/luvx/upx:latest --link /usr/bin/upx /usr/bin/
WORKDIR /build

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux
COPY . .
RUN go mod download && go mod tidy && go build -ldflags="-w -s" -v -o app main.go && upx -9 app

# -------------------------------------------

FROM ${REGISTRY}alpine:3.22

WORKDIR /app
ENV log_LogDir=/root/data/logs TZ=Asia/Shanghai
COPY --from=builder /build/app .

# 将时区设置为东八区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo "Asia/Shanghai" > /etc/timezone # && apk del tzdata

EXPOSE 8090

CMD ["./app"]
