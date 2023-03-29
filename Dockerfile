FROM golang:1.20-alpine

# 安装必要的工具和依赖
RUN apk update && apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 拷贝代码到容器中
COPY . .

# 下载依赖
RUN go get -d -v

# 编译程序
RUN go build -o app

RUN go build -o app

# 暴露端口
EXPOSE 80

# 运行程序
CMD ["./app"]