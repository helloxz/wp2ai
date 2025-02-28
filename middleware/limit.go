package middleware

import (
	"fmt"
	"strconv"
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func LimitChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取访客的IP
		ip := utils.GetClientIP(c)
		// 设置一个key
		key := fmt.Sprintf("limit_chat:%s", ip)
		// key转为[]byte
		keyByte := []byte(key)
		// 从内存中获取
		value, _ := utils.Cache.Get(keyByte)
		// 如果为空
		if value == nil {
			// 设置一个值
			utils.Cache.Set(keyByte, []byte("0"), 60*60*24)
			// 继续执行后续的逻辑
			c.Next()
			return
		}
		// 如果不为空，则把value转为int
		valueStr := string(value)
		count, _ := strconv.Atoi(valueStr)
		// 获取配置中的限制
		limitNum := viper.GetInt("app.chat_limit")
		// 如果超过限制
		if count >= limitNum {
			// 返回错误信息
			c.JSON(200, gin.H{
				"code": 429,
				"msg":  "Too many requests",
				"data": "",
			})
			c.Abort()
			return
		}
		// 如果没有超过限制，则继续执行后续的逻辑
		c.Next()
	}
}
