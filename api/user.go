package api

import (
	"strings"
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 初始化用户名
func InitUser(c *gin.Context) {
	// 通过POST获取用户名
	email := c.PostForm("email")
	// 转为小写
	email = strings.ToLower(email)
	// 通过POST获取密码
	password := c.PostForm("password")

	// 获取配置文件中的邮箱和密码
	configEmail := viper.GetString("auth.email")
	configPassword := viper.GetString("auth.password")

	// 如果邮箱和密码其中一个不为空，则不允许初始化
	if configEmail != "" || configPassword != "" {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "The system has been initialized",
			"data": "",
		})
		c.Abort()
		return
	}

	// 验证邮箱是否有效
	if !utils.IsEmail(email) {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Invalid email address",
			"data": "",
		})
		c.Abort()
		return
	}

	// 验证密码是否有效
	if !utils.IsPassword(password) {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Invalid password",
			"data": "",
		})
		c.Abort()
		return
	}

	passTxt := email + password + "wp2ai"
	// 生成MD5哈希
	passHash := utils.MD5(passTxt)

	// 写入配置文件
	if !utils.SetConfigStr("auth.email", email) {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Failed to write configuration file",
			"data": "",
		})
		c.Abort()
		return
	}
	// 写入密码
	if !utils.SetConfigStr("auth.password", passHash) {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Failed to write configuration file",
			"data": "",
		})
		c.Abort()
		return
	}
	// 提示成功
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "",
	})
}

// 用户登录
func Login(c *gin.Context) {
	// 通过POST获取用户名
	email := c.PostForm("email")
	// 转为小写
	email = strings.ToLower(email)
	// 通过POST获取密码
	password := c.PostForm("password")

	// 获取配置文件中的邮箱和密码
	configEmail := viper.GetString("auth.email")
	configPassword := viper.GetString("auth.password")

	// 如果配置文件中的邮箱或者密码任意一个为空，则提示应该初始化
	if configEmail == "" || configPassword == "" {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Please initialize the system",
			"data": "",
		})
		c.Abort()
		return
	}

	// 验证邮箱是否有效
	if !utils.IsEmail(email) {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Invalid email address",
			"data": "",
		})
		c.Abort()
		return
	}

	// 验证密码是否有效
	if !utils.IsPassword(password) {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Invalid password",
			"data": "",
		})
		c.Abort()
		return
	}

	// 验证邮箱和密码是否正确
	if email != configEmail {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Username or password is incorrect",
			"data": "",
		})
		c.Abort()
		return
	}

	passTxt := email + password + "wp2ai"
	// 生成MD5哈希
	passHash := utils.MD5(passTxt)

	if passHash != configPassword {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Username or password is incorrect",
			"data": "",
		})
		c.Abort()
		return
	}

	// 生成一个token
	token := "web-" + utils.RandStr(28)
	// 取出前8位作为key
	key := token[:8]
	keyByte := []byte(key)
	valueByte := []byte(token)

	// 设置过期时间为 30 天
	expireTime := 30 * 24 * 60 * 60 // 30 天的秒数
	// 存储到缓存中
	err := utils.Cache.Set(keyByte, valueByte, expireTime)

	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to write cache",
			"data": "",
		})
		c.Abort()
		return
	}

	// 提示成功
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": token,
	})
}

// 检测用户是否登录
func IsLogin(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "",
	})
}
