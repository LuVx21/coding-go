ARG CR=${CR:-""}
ARG REGISTRY=${CR:+$CR/}
FROM ${REGISTRY}golang:1.26-alpine3.23 AS builder

COPY --from=registry.cn-shanghai.aliyuncs.com/luvx/upx:latest --link /usr/bin/upx /usr/bin/
WORKDIR /build

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux
RUN go install -ldflags="-w -s" github.com/luvx21/coding-go/cmd/kv@latest && upx -9 /go/bin/kv
COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -v -o app main.go && upx -9 app

# -------------------------------------------

FROM ${REGISTRY}alpine:3.23

RUN apk add --no-cache libgcc gcompat
WORKDIR /app
ENV log_LogDir=/root/data/logs TZ=Asia/Shanghai
COPY --from=builder /build/app .
COPY --from=builder /go/bin/kv /usr/bin/

# 将时区设置为东八区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo "Asia/Shanghai" > /etc/timezone # && apk del tzdata

EXPOSE 8090

CMD ["./app"]
