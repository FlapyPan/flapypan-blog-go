package conf

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// LoadEnv 加载 .env 文件
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("加载 .env 文件失败 %e", err)
	}
}

// Env 读取配置
func Env(key string) string {
	return os.Getenv(key)
}
