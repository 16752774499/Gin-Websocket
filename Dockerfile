# 使用 golang 作为基础镜像
FROM golang:1.23.5 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制当前目录内容到工作目录
COPY . .

# 构建 Go 程序
RUN go build -o main .

# 使用一个轻量级的镜像来减小最终镜像大小
FROM alpine:latest

# 将构建好的程序从 builder 复制过来
WORKDIR /root/
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 3333

# 运行程序
CMD ["./main"]
