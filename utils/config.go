package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"wp2ai/model"

	"github.com/spf13/viper"
)

var (
	once sync.Once
)

// 全局就一个init函数，避免其它地方再次声明，以免逻辑容易出错
func InitConfig() {
	once.Do(func() {
		// 创建必要的目录
		CreateDir("data/config")
		CreateDir("data/db")
		CreateDir("data/logs")

		//默认配置文件
		config_file := "data/config/config.toml"
		// 检查配置文件是否存在，不存在则复制一份
		if _, err := os.Stat(config_file); os.IsNotExist(err) {
			// 创建目录
			os.MkdirAll("data/config", os.ModePerm)
			// 复制配置文件
			err := CopyFile("config.toml", config_file)
			if err != nil {
				fmt.Println("Failed to copy config file:", err)
				os.Exit(1)
			}
		}

		viper.SetConfigFile(config_file) // 指定配置文件路径
		//指定ini类型的文件
		viper.SetConfigType("toml")
		err := viper.ReadInConfig() // 读取配置信息
		if err != nil {             // 读取配置信息失败
			// 写入日志
			fmt.Println("Failed to read config:", err)
			os.Exit(1)
		}
	})

	// 连接内存缓存
	InitCache()

	// 连接数据库
	model.InitDB()
	model.InitVecDB()
	// Test()
	go model.InitWPDB()
	// 异步执行定时任务
	go InitCrontab()

	// 初始化密钥
	go InitToken()
}

// 初始化密钥
func InitToken() {
	// 获取密钥
	token := viper.GetString("auth.token")
	// 如果密钥为空
	if token == "" {
		tokenStr := "sk-" + RandString(29)
		// 设置密钥
		viper.Set("auth.token", tokenStr)
		// 写入配置
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("Failed to init token:", err)
		}
	}
}

// 生成一个随机字符串
func RandString(length int) string {
	// 定义字符集
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	// 将字符集转换为 rune 切片，以便随机选择字符
	charsetRunes := []rune(charset)

	// 创建一个 strings.Builder，用于高效构建字符串
	var sb strings.Builder
	// 设置 strings.Builder 的初始容量，避免频繁扩容
	sb.Grow(length)

	// 使用当前时间作为随机数种子，确保每次运行生成不同的随机字符串
	rand.Seed(time.Now().UnixNano())

	// 循环生成随机字符，并添加到 strings.Builder 中
	for i := 0; i < length; i++ {
		// 从字符集中随机选择一个字符
		randomIndex := rand.Intn(len(charsetRunes))
		randomChar := charsetRunes[randomIndex]
		// 将随机字符添加到 strings.Builder 中
		sb.WriteRune(randomChar)
	}

	// 返回生成的随机字符串
	return sb.String()
}

// copyFile 复制文件内容
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// 接收一个路径作为参数，判断路径是否存在，不存在则创建
func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("无法创建目录：%w", err)
		}
	}
	return nil
}
