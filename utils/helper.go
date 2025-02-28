package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// 设置配置文件
func SetConfigStr(key string, value string) bool {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		// 写入错误日志
		WriteLog(fmt.Sprintf("%s", err))
		return false
	} else {
		return true
	}
}

// 生成一个随机字符串
func RandStr(length int) string {
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

// 获取客户端IP
func GetClientIP(c *gin.Context) string {
	//尝试通过X-Forward-For获取
	ip := c.Request.Header.Get("X-Forward-For")
	//如果没获取到，则通过X-real-ip获取
	if ip == "" {
		ip = c.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		//依然没获取到，则通过gin自身方法获取
		ip = c.ClientIP()
	}
	//判断IP格式是否正确，避免伪造IP
	if net.ParseIP(ip) == nil {
		ip = "0.0.0.0"
	}
	return ip
}
