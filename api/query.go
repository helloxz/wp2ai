package api

import (
	"fmt"
	"wp2ai/model"
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
)

func Query(c *gin.Context) {
	// 通过GET获取keywords
	keywords := c.Query("keywords")

	// 向量化
	input, err := utils.DataEmbedding(keywords)
	if err != nil {
		utils.WriteLog(err.Error())
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "Data embedding failed",
			"data": "",
		})
		c.Abort()
		return
	}
	// 查询数据
	result, err := model.GetDocument(input)
	fmt.Println(err)
	if err != nil {
		utils.WriteLog(err.Error())
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "Query failed",
			"data": "",
		})
		c.Abort()
		return
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": result,
	})
}
