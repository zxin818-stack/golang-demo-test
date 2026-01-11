FROM golang:latest

# 创建/config_file目录并设置777权限（所有用户可读写执行）
RUN mkdir -p /config_file && chmod 777 /config_file

# 复制本地编译好的可执行文件
COPY golang-demo-test /main

# 复制配置文件
#COPY config.yaml /config.yaml

# 暴露端口
EXPOSE 8080

# 运行程序
CMD ["/main"]