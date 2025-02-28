package api

import (
	"strconv"
	"wp2ai/model"

	"github.com/gin-gonic/gin"
)

// 添加单个数据
func AddPost(c *gin.Context) {
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

	// 获取文章
	wpPost, _ := model.GetWpPostById(idUint)
	// 如果post为空
	if wpPost == nil {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "Post not found",
			"data": "",
		})
		c.Abort()
		return
	}

	var post model.Post
	post.PostID = wpPost.ID
	post.PostTitle = wpPost.PostTitle
	post.PostDate = wpPost.PostDate
	post.PostContent = wpPost.PostContent
	post.Status = 0

	// 插入数据库
	err = model.InsertPost(post)
	// 如果出错
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Insert failed",
			"data": "",
		})
		c.Abort()
		return
	}
	// 成功了
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "",
	})
}
