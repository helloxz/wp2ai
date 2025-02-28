package api

import (
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 需要授权，给管理员使用
func AppInfo(c *gin.Context) {
	var data = map[string]string{
		"app.doc_limit":         viper.GetString("app.doc_limit"),
		"app.chat_limit":        viper.GetString("app.chat_limit"),
		"wordpress.domain":      viper.GetString("wordpress.domain"),
		"embedding.key":         viper.GetString("embedding.key"),
		"embedding.url":         viper.GetString("embedding.url"),
		"embedding.model":       viper.GetString("embedding.model"),
		"openai.key":            viper.GetString("openai.key"),
		"openai.url":            viper.GetString("openai.url"),
		"openai.model":          viper.GetString("openai.model"),
		"wordpress.db_host":     viper.GetString("wordpress.db_host"),
		"wordpress.db_name":     viper.GetString("wordpress.db_name"),
		"wordpress.db_password": viper.GetString("wordpress.db_password"),
		"wordpress.db_username": viper.GetString("wordpress.db_username"),
		"version":               utils.Version,
		"email":                 viper.GetString("auth.email"),
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

// 无需授权，给访客使用
func SiteInfo(c *gin.Context) {
	isInit := "yes"
	// 检查系统是否初始化
	email := viper.GetString("auth.email")
	password := viper.GetString("auth.password")
	// 如果两个都为空，则系统未初始化
	if email == "" && password == "" {
		isInit = "no"
	}
	var data = map[string]string{
		"wp_domain": viper.GetString("wordpress.domain"),
		"is_init":   isInit,
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}
