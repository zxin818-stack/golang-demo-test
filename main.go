package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

// Config 结构体定义配置项
type Config struct {
	AppName  string         `json:"app_name"`
	Version  string         `json:"version"`
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Features []string       `json:"features"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Port    int `json:"port"`
	Timeout int `json:"timeout"`
}

func main() {
	// 从环境变量获取配置文件路径
	configPath := os.Getenv("LOCAL_CONFIG_PATH")
	if configPath == "" {
		// 如果环境变量未设置，使用默认路径
		configPath = "config.yaml"
		fmt.Printf("环境变量 LOCAL_CONFIG_PATH 未设置，使用默认路径: %s\n", configPath)
	} else {
		fmt.Printf("使用配置文件路径: %s\n", configPath)
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析YAML配置
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 打印配置信息
	fmt.Println("=== 配置信息 ===")
	fmt.Printf("应用名称: %s\n", config.AppName)
	fmt.Printf("版本: %s\n", config.Version)
	fmt.Println("\n数据库配置:")
	fmt.Printf("  主机: %s\n", config.Database.Host)
	fmt.Printf("  端口: %d\n", config.Database.Port)
	fmt.Printf("  用户名: %s\n", config.Database.Username)
	fmt.Printf("  密码: %s\n", config.Database.Password)
	fmt.Println("\n服务器配置:")
	fmt.Printf("  端口: %d\n", config.Server.Port)
	fmt.Printf("  超时时间: %d秒\n", config.Server.Timeout)
	fmt.Println("\n功能列表:")
	for i, feature := range config.Features {
		fmt.Printf("  %d. %s\n", i+1, feature)
	}

	// 启动HTTP服务器
	startHTTPServer(&config)
}

// getConfigHandler 处理 /getconfig 请求
func getConfigHandler(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头为JSON格式
		w.Header().Set("Content-Type", "application/json")

		// 将配置转换为JSON格式
		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			http.Error(w, "Failed to marshal config to JSON", http.StatusInternalServerError)
			return
		}

		// 返回配置信息
		w.WriteHeader(http.StatusOK)
		w.Write(configJSON)
	}
}

// startHTTPServer 启动HTTP服务器
func startHTTPServer(config *Config) {
	// 设置路由
	http.HandleFunc("/getconfig", getConfigHandler(config))

	// 获取服务器端口，如果配置中未设置则使用默认端口8080
	serverPort := config.Server.Port
	if serverPort == 0 {
		serverPort = 8080
	}

	fmt.Printf("\n=== HTTP服务器启动 ===\n")
	fmt.Printf("服务器运行在端口: %d\n", serverPort)
	fmt.Printf("访问 http://localhost:%d/getconfig 获取配置信息\n", serverPort)

	// 启动服务器
	err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
	if err != nil {
		log.Fatalf("HTTP服务器启动失败: %v", err)
	}
}
