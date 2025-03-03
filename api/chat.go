package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"wp2ai/model"
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

// json结构体
type Input struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

var ChatClient *openai.Client

func Chat(c *gin.Context) {
	var inputs []Input
	// 绑定json数据
	if err := c.ShouldBindJSON(&inputs); err != nil {
		// 返回错误信息
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Incorrect format!",
			"data": "",
		})
		c.Abort()
		return
	}
	var extMsg []openai.ChatCompletionMessage
	// 遍历inputs
	for _, input := range inputs {
		// 判断type是否为user
		if input.Type == "user" {
			// 添加到extMsg中
			extMsg = append(extMsg, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: input.Msg,
			})
		} else if input.Type == "ai" {
			// 添加到extMsg中
			extMsg = append(extMsg, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: input.Msg,
			})

		}
	}
	// 取出最后一个元素(用户的输入)
	input := extMsg[len(extMsg)-1].Content
	// 删除extMsg最后一个元素
	// extMsg = extMsg[:len(extMsg)-1]
	// 通过Post获取用户输入
	// input := c.PostForm("input")
	// fmt.Println(input)

	modelName := viper.GetString("openai.model")
	// 初始化客户端
	InitChatClient()

	// 将用户最后输入的内容向量化
	embedding, err := utils.DataEmbedding(input)
	// fmt.Println(embedding)
	// 如果出现错误
	if err != nil {
		// 返回错误信息
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Data embedding failed",
			"data": "",
		})
		// 终止
		c.Abort()
		return
	}
	// 查询出匹配的数据
	results, err := model.GetDocument(embedding)

	// 写入日志
	utils.WriteLog(fmt.Sprintf("GetDocument error: %v", err))
	// 如果出现错误
	if err != nil {
		// 返回错误信息
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Get document failed",
			"data": "",
		})
		// 终止
		c.Abort()
		return
	}
	var messages []openai.ChatCompletionMessage
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "你是一名AI助手，名字叫WP2AI，以下是一些资料，你需要先理解并掌握。然后根据掌握的资料来回答问题，如果用户的问题不在资料中，你需要拒绝回答。",
	})
	// 引用结果
	quotes := "  \n  \n> **引用内容：**  \n "
	// fmt.Println(results)
	// 遍历结果
	allDistancesLarge := true // 标记是否所有 Distance 都大于等于 1.0

	// 检查是否所有 Distance 都大于等于 1.0
	for _, result := range results {
		if result.Distance < 1.0 {
			allDistancesLarge = false
			break
		}
	}
	// 部分向量模型数值原因，暂时停用这个功能
	allDistancesLarge = false

	// 如果所有allDistancesLarge都大于1.0，则去掉results最后2行
	if allDistancesLarge {
		results = results[:len(results)-2]
	}

	domain := viper.GetString("wordpress.domain")

	for index, result := range results {
		// fmt.Println(result.Distance)
		// 添加前缀
		prefix := fmt.Sprintf("资料%d：", index+1)
		// 添加到messages中
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: prefix + result.Content,
		})
		quotes += fmt.Sprintf("> %v. [%v](%v/?p=%v)  \n  ", index+1, result.Title, domain, result.PostID)
	}
	// 组合消息，已经学习的内容 + 前端传递的内容
	// messages = append(messages, extMsg...)
	// 最后一次用户输入的内容
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	// fmt.Println(messages)

	// 设置响应header头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   true,
	}
	// fmt.Println(messages)
	stream, err := ChatClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Println(err)
		c.SSEvent("data", toJsonStr(fmt.Sprintf("ChatCompletionStream error: %v", err)))
		c.Writer.Flush() // 添加这一行来刷新缓冲区
		c.SSEvent("close", "[DONE]")
		c.Writer.Flush() // 添加这一行来刷新缓冲区
		c.Abort()
		return
	}

	// fmt.Println("dsdsd")
	// 输出消息
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			c.SSEvent("data", toJsonStr(quotes))
			c.Writer.Flush() // 添加这一行来刷新缓冲区
			c.SSEvent("close", "[DONE]")
			c.Writer.Flush() // 添加这一行来刷新缓冲区
			// 请求次数+1
			ChatCount(c)
			return
		}

		if err != nil {
			c.SSEvent("data", toJsonStr(fmt.Sprintf("Stream error: %v", err)))
			c.Writer.Flush() // 添加这一行来刷新缓冲区
			c.SSEvent("close", "[DONE]")
			c.Writer.Flush() // 添加这一行来刷新缓冲区
			c.Abort()
			return
		}

		// 获得接口返回的内容
		content := response.Choices[0].Delta.Content
		// 追加流消息
		// data += content
		c.SSEvent("data", toJsonStr(content))
		c.Writer.Flush() // 添加这一行来刷新缓冲区
	}
	// 可选：发送完成标志
	// c.SSEvent("close", "[DONE]")
	// c.Writer.Flush() // 添加这一行来刷新缓冲区

}

// 初始化AI客户端
func InitChatClient() {
	config := openai.DefaultConfig(viper.GetString("openai.key"))
	apiURL := viper.GetString("openai.url")
	// modelName := viper.GetString("openai.model")
	config.BaseURL = apiURL
	// 判断ChatClient是否为空
	if ChatClient == nil {
		// 初始化ChatClient
		ChatClient = openai.NewClientWithConfig(config)
		// 初始化成功提示
		fmt.Printf("ChatClient connection succeeded!\n")
	}
}

// 将字符串转为json字符串输出
func toJsonStr(input string) string {
	data := struct {
		V string `json:"v"`
	}{
		V: input,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

// 让对话次数+1
func ChatCount(c *gin.Context) {
	// 获取访客的IP
	ip := utils.GetClientIP(c)
	// 设置一个key
	key := fmt.Sprintf("limit_chat:%s", ip)
	// key转为[]byte
	keyByte := []byte(key)
	// 从内存中获取
	value, _ := utils.Cache.Get(keyByte)
	// 如果为空
	if value == nil {
		return
	}
	// 如果不为空，则把value转为int
	valueStr := string(value)
	count, _ := strconv.Atoi(valueStr)
	// 打印次数
	// fmt.Println(count)
	// 次数+1
	count++
	// 设置一个值
	err := utils.Cache.Set(keyByte, []byte(strconv.Itoa(count)), 60*60*24)
	if err != nil {
		// 写入日志
		utils.WriteLog(fmt.Sprintf("%s", err))
	}
}
