package api

import (
	"fmt"
	"wp2ai/model"
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Hello World",
		"data": "",
	})
}

// 批量扫描WordPress中的文章ID，然后入库
func BatchScan(c *gin.Context) {
	// 再次初始化一次
	model.InitWPDB()
	// 获取现有行数
	count := model.CountPosts(-1)
	// 如果行数大于0，说明已经扫描过了，不需要再扫描
	if count > 0 {
		c.JSON(200, gin.H{
			"code": -1000,
			"msg":  "已经扫描过了，请勿重复扫描！",
			"data": "",
		})
		c.Abort()
		return
	}
	postIds := model.GetPostIds()

	var posts []model.Post
	var post model.Post
	// 遍历postIds，把所有id全部给posts
	for _, postId := range postIds {
		post.PostID = postId.ID
		posts = append(posts, post)
	}
	// 插入到数据库中
	err := model.InsertPosts(posts)
	if err != nil {
		utils.WriteLog(fmt.Sprintf("批量插入文章ID失败：%v", err))
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Scan failed",
			"data": "",
		})
		c.Abort()
		return
	}
	// 返回结果
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": len(posts),
	})
}
