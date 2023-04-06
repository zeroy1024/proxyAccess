# 第一阶段：编译二进制文件
FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o ProxyAccess .

# 第二阶段：构建最终镜像
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/ProxyAccess .
CMD ["/app/ProxyAccess"]
