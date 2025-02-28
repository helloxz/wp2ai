package middleware

import (
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 写一个认证中间件，从header的Authorization中获取Bearer token，然后解析token，如果token合法，则继续执行后续的逻辑，否则返回401
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header中获取Authorization
		bearer := c.GetHeader("Authorization")
		// 如果token为空
		if bearer == "" {
			// 返回401
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
				"data": "",
			})
			c.Abort()
			return
		}
		// 如果长度小于7
		if len(bearer) <= 7 {
			// 返回401
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
				"data": "",
			})
			c.Abort()
			return
		}
		// 从bearer中获取token
		token := bearer[7:]
		// 打印token
		// fmt.Println(token)
		// 如果token 为空
		if token == "" {
			// 返回401
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
				"data": "",
			})
			c.Abort()
			return
		}
		// 获取配置文件里面的token
		configToken := viper.GetString("auth.token")
		if token == configToken {
			// 继续执行后续的逻辑
			c.Next()
			return
		}
		// 这里还在执行，说明与配置文件里面的token不匹配，则需要从内存获取
		key := token[:8]
		keyByte := []byte(key)
		// 通过内存获取
		valueByte, err := utils.Cache.Get(keyByte)
		// 如果获取失败
		if err != nil {
			// 返回401
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
				"data": "",
			})
			c.Abort()
			return
		}
		// 将valueByte转为字符串
		value := string(valueByte)
		// 如果token不匹配
		if token != value {
			// 返回401
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "Unauthorized",
				"data": "",
			})
			c.Abort()
			return
		}

		// 上面全部通过了，继续执行后续的逻辑
		c.Next()
	}
}
