package api

import (
	"wp2ai/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 设置配置
func SetConfig(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Incorrect format!",
			"data": "",
		})
		c.Abort()
		return
	}

	// 遍历 JSON 数据
	for key, value := range data {
		// 设置配置
		viper.Set(key, value)
		// 写入配置
		err := viper.WriteConfig()
		if err != nil {
			c.JSON(200, gin.H{
				"code": 500,
				"msg":  "Failed to write config!",
				"data": "",
			})
			c.Abort()
			return
		}

		// 如果key == db_name
		if key == "wordpress.db_name" {
			// 清空WP
			model.WP = nil
			// 重新初始化
			go model.InitWPDB()
		}

		// 如果key ==openai.key
		if key == "openai.key" {
			// 清空ChatClient
			ChatClient = nil
			// 重新初始化
			go InitChatClient()
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "",
	})
}
