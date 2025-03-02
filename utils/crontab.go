package utils

import (
	"fmt"
	"strconv"
	"time"
	"wp2ai/model"

	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func InitCrontab() {
	// 创建一个定时任务的实例，带上cron.WithSeconds()则精确到秒级别
	c := cron.New()
	// 添加定时任务,每隔5分钟探测一次所有节点
	c.AddFunc("*/1 * * * *", Vectorization)

	// 启动定时任务
	c.Start()

	// 这里使用select{}可以使主程序持续运行，否则主程序可能会结束，导致定时任务未能执行
	select {}
}

// 每次获取10条数据
func Vectorization() {
	minProcess := viper.GetInt("app.min_process")
	// fmt.Println("minProcess:", minProcess)
	// 获取文章数据
	posts := model.GetPosts(minProcess)
	// 如果没有数据，则直接返回
	if len(posts) == 0 {
		return
	}
	// 遍历posts，得到文章IDS
	var ids []uint
	for _, post := range posts {
		ids = append(ids, post.PostID)
	}
	// 批量更新这批文章的状态为1：处理中
	err := model.UpdatePostsStatus(ids, 1)
	if err != nil {
		// 写入日志
		WriteLog(fmt.Sprintf("批量更新文章状态失败：%v", err))
		return
	}
	// 查询出这些文章
	wpPosts := model.GetPostsByIds(ids)
	// 遍历wpPosts，然后做向量化处理
	for _, wpPost := range wpPosts {
		go SingleVectorization(wpPost)
		time.Sleep(time.Millisecond * 500) // 休眠 0.5 秒，避免并发过高
	}
}

// 单个处理数据
func SingleVectorization(wpPost model.WpPost) {
	var post model.Post
	post.PostID = wpPost.ID
	post.Status = 1
	// 如果文章标题或者内容是空的，则跳过不要处理
	if wpPost.PostTitle == "" || wpPost.PostContent == "" {
		return
	}
	// 更新文章状态为1：处理中
	err := model.UpdatePost(post)
	if err != nil {
		// 写入日志
		WriteLog(fmt.Sprintf("post_id:%v更新文章状态失败：%v", post.PostID, err))
		return
	}
	// 获取文章内容
	content := wpPost.PostContent
	// 获取向量化数据
	embedding, err := DataEmbedding(content)
	if err != nil {
		// 写入日志
		WriteLog(fmt.Sprintf("post_id:%v,文章向量化失败：%v", post.PostID, err))
		// 更新文章状态为4：存在错误
		post.Status = 4
		// 更新文章
		model.UpdatePost(post)
		return
	}
	// 更新文章状态为3：已处理
	post.Status = 3
	post.PostTitle = wpPost.PostTitle
	// post.PostDate = wpPost.PostDate
	// post.UpdatedAt = &wpPost.PostDate
	post.PostContent = content
	post.PostDate = wpPost.PostDate

	err = model.UpdatePost(post)
	if err != nil {
		// 写入日志
		WriteLog(fmt.Sprintf("post_id:%v更新文章状态失败：%v", post.PostID, err))
		return
	}
	// 插入到items表中
	item := model.Item{
		PostID:    strconv.FormatUint(uint64(post.PostID), 10),
		Embedding: embedding,
		Title:     wpPost.PostTitle,
		Content:   content,
	}
	// 插入到vector数据库中
	err = model.InsertDocument(item)
	if err != nil {
		// 写入日志
		WriteLog(fmt.Sprintf("post_id:%v,插入向量数据失败：%v", post.PostID, err))
		// 更新文章状态为4：存在错误
		post.Status = 4
		// 更新文章
		model.UpdatePost(post)
		return
	}
}

// 数据向量化
func DataEmbedding(input string) ([]float32, error) {
	client := resty.New()

	// 请求的URL和headers
	url := viper.GetString("embedding.url") // 从配置文件中获取
	headers := map[string]string{
		"Authorization": "Bearer " + viper.GetString("embedding.key"), // 替换为你的 API Key
		"Content-Type":  "application/json",
	}

	// 请求的body
	body := map[string]interface{}{
		"input":           input,
		"model":           viper.GetString("embedding.model"), // 替换为你的模型名称
		"encoding_format": "float",
	}

	// 发送POST请求
	resp, err := client.R().
		SetHeaders(headers).
		SetBody(body).
		Post(url)

	if err != nil {
		// 写入日志
		WriteLog(fmt.Sprintf("请求向量化接口失败：%v", err))
		return nil, err
	}

	// 获取响应数据
	resBody := resp.Body()
	// 转为字符串
	bodyStr := string(resBody)
	// 获取 "embedding" 数组
	embeddingStr := gjson.Get(bodyStr, "data.0.embedding")

	// 提取嵌套数组并将其转换为 []float32
	var embedding []float32
	for _, v := range embeddingStr.Array() {
		// 将字符串转为 float32 并添加到切片中
		f, err := strconv.ParseFloat(v.String(), 32)
		if err != nil {
			// log.Fatalf("Error parsing float: %v", err)
			// 写入日志
			WriteLog(fmt.Sprintf("Error parsing float: %v", err))
		}
		embedding = append(embedding, float32(f))
	}

	return embedding, nil

}
