package api

import (
	"strconv"
	"wp2ai/model"

	"github.com/gin-gonic/gin"
)

// 清空所有数据
func DeleteAll(c *gin.Context) {
	// 清空posts表
	err := model.TruncatePosts()
	// 如果出错
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Delete posts failed",
			"data": "",
		})
		c.Abort()
		return
	}
	// 继续清空
	err = model.TruncateItems()
	// 如果出错
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Delete embedding failed",
			"data": "",
		})
		c.Abort()
		return
	}

	// 返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "",
	})

}

// 删除单个数据
func DeletePost(c *gin.Context) {
	// 通过POST获取id
	id := c.PostForm("id")
	// 将id转为utint
	idUint64, err := strconv.ParseUint(id, 10, 0)
	idUint := uint(idUint64)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Invalid id",
			"data": "",
		})
		c.Abort()
		return
	}

	err = model.DeletePost(idUint)
	// 如果出错
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Delete failed",
			"data": "",
		})
		c.Abort()
		return
	}

	// 继续删除，只管删除，不管是否成功
	err = model.DeleteItem(id)

	// 返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "",
	})

}
