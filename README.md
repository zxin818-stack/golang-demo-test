# Go配置读取Demo

这是一个简单的Go语言demo，演示如何从本地配置文件中读取配置项，配置文件路径通过环境变量`LOCAL_CONFIG_PATH`获取。

## 项目结构

```
golang-demo-test/
├── go.mod          # Go模块文件
├── main.go         # 主程序文件
├── config.yaml     # 示例配置文件
└── README.md       # 项目说明文档
```

## 功能特性

- 从环境变量`LOCAL_CONFIG_PATH`读取配置文件路径
- 支持YAML格式配置文件
- 如果环境变量未设置，默认使用当前目录下的`config.yaml`文件
- 解析并打印配置信息
- 提供HTTP服务器，支持`/getconfig`接口返回JSON格式配置数据

## 使用方法

### 1. 设置环境变量（可选）

```bash
# 设置自定义配置文件路径
export LOCAL_CONFIG_PATH="/path/to/your/config.json"

# 在Windows上使用：
set LOCAL_CONFIG_PATH="C:\path\to\your\config.json"
```

### 2. 运行程序

```bash
# 如果未设置环境变量，程序会自动使用默认的config.yaml文件
go run main.go

# 如果设置了环境变量，程序会使用指定路径的配置文件
export LOCAL_CONFIG_PATH="/path/to/config.yaml"
go run main.go
```

### 3. 预期输出

程序会输出类似以下内容：

```
使用配置文件路径: config.yaml
=== 配置信息 ===
应用名称: golang-demo-test
版本: 1.0.0

数据库配置:
  主机: localhost
  端口: 5432
  用户名: admin
  密码: secret

服务器配置:
  端口: 8080
  超时时间: 30秒

功能列表:
  1. auth
  2. logging
  3. caching

=== HTTP服务器启动 ===
服务器运行在端口: 8080
访问 http://localhost:8080/getconfig 获取配置信息
```

## HTTP服务器功能

程序启动后会同时启动一个HTTP服务器，提供以下接口：

### /getconfig 接口

- **URL**: `http://localhost:8080/getconfig` (端口号根据配置文件中的server.port设置)
- **方法**: GET
- **响应格式**: JSON
- **功能**: 返回所有配置参数的JSON格式数据

### 示例请求

```bash
# 使用curl获取配置信息
curl http://localhost:8080/getconfig
```

### 示例响应

```json
{
  "app_name": "golang-demo-test",
  "version": "1.0.0",
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "admin",
    "password": "secret"
  },
  "server": {
    "port": 8080,
    "timeout": 30
  },
  "features": [
    "auth",
    "logging",
    "caching"
  ]
}
```

### 注意事项

- HTTP服务器端口使用配置文件中的`server.port`设置，如果未设置则默认使用8080端口
- 程序启动后会同时打印配置信息到控制台和启动HTTP服务器
- 可以通过配置文件灵活调整服务器端口和其他参数

## Docker部署

### 构建流程

1. **首先在本地编译Go程序**：
```bash
# 使用静态链接编译，确保在scratch镜像中正常运行
CGO_ENABLED=0 GOOS=linux go build -o main .
```

2. **构建Docker镜像**：
```bash
# 在项目根目录下构建Docker镜像（使用scratch镜像，无需网络依赖）
docker build -t golang-demo-test .

# 或者指定标签版本
docker build -t golang-demo-test:latest .
```

**注意**：Dockerfile使用多阶段构建：
- 第一阶段使用Alpine镜像创建/config_file目录并设置777权限（所有用户可读写执行）
- 第二阶段使用scratch镜像运行程序
- 请确保在构建Docker镜像前已经执行了本地编译步骤
- 构建过程需要网络连接来下载Alpine镜像

### 使用构建脚本（推荐）

项目提供了一个构建脚本 `build.sh` 来自动化整个构建流程：

```bash
# 给脚本添加执行权限
chmod +x build.sh

# 运行构建脚本
./build.sh
```

该脚本会自动完成以下步骤：
1. 编译Go程序（静态链接）
2. 构建Docker镜像
3. 提供运行容器的命令示例

### 运行Docker容器

```bash
# 基本运行（使用镜像内的默认配置文件）
docker run -p 8080:8080 golang-demo-test

# 使用自定义配置文件（挂载本地配置文件）
docker run -p 8080:8080 -v /path/to/your/config.yaml:/root/config.yaml golang-demo-test

# 后台运行
docker run -d -p 8080:8080 --name config-server golang-demo-test

# 使用自定义环境变量
docker run -p 8080:8080 -e LOCAL_CONFIG_PATH=/custom/path/config.yaml golang-demo-test
```

### 验证部署

```bash
# 检查容器是否运行
docker ps

# 查看容器日志
docker logs config-server

# 测试HTTP接口
curl http://localhost:8080/getconfig
```

### Docker Compose部署（可选）

创建`docker-compose.yml`文件：

```yaml
version: '3.8'
services:
  config-server:
    image: golang-demo-test:latest
    build: .
    ports:
      - "8080:8080"
    environment:
      - LOCAL_CONFIG_PATH=/root/config.yaml
    volumes:
      - ./config.yaml:/root/config.yaml
```

运行：
```bash
docker-compose up -d
```

## 配置文件格式

配置文件使用YAML格式，支持嵌套结构：

```yaml
app_name: 应用名称
version: 版本号

database:
  host: 数据库主机
  port: 5432
  username: 用户名
  password: 密码

server:
  port: 8080
  timeout: 30

features:
  - 功能1
  - 功能2
  - 功能3
```

## 自定义配置

您可以修改`config.json`文件来添加自己的配置项，或者在`main.go`中扩展`Config`结构体来支持新的配置字段。