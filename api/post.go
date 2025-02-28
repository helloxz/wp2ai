package api

import (
	"strconv"
	"wp2ai/model"

	"github.com/gin-gonic/gin"
)

// 声明一个结构体，用于返回信息
type RePost struct {
	Posts []model.Post `json:"posts"`
	Count int          `json:"count"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}

// 显示所有文章
func PostList(c *gin.Context) {
	// 获取页码
	page := c.DefaultQuery("page", "1")
	// 获取每页数量
	limit := c.DefaultQuery("limit", "10")
	// 转为int
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	// fmt.Println(pageInt, limitInt)
	var rePost RePost
	// 查询数据
	posts, err := model.GetPostsByPage(pageInt, limitInt)

	// fmt.Println("dsd", posts)

	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "Query failed",
			"data": "",
		})
		c.Abort()
		return
	}

	rePost.Posts = posts
	rePost.Count = model.CountPosts(-1)
	rePost.Page = pageInt
	rePost.Limit = limitInt

	// 返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": rePost,
	})
}
