package router

import (
	"wp2ai/api"
	"wp2ai/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Start() {
	//gin运行模式
	RunMode := viper.GetString("server.mode")
	// 设置运行模式
	gin.SetMode(RunMode)
	//运行gin
	r := gin.Default()
	// 全局跨域中间件
	r.Use(middleware.CORSMiddleware())

	//前台首页
	r.LoadHTMLFiles("assets/templates/index.html")
	r.GET("/", api.Home)
	r.GET("/index.html", api.Home)
	r.Static("/assets", "assets")
	r.NoRoute(api.Home)

	r.GET("/api/batch-scan", middleware.Auth(), api.BatchScan)
	r.GET("/api/query", api.Query)
	r.POST("/api/chat", middleware.LimitChat(), api.Chat)
	r.GET("/api/post/list", middleware.Auth(), api.PostList)
	r.GET("/api/get/appinfo", middleware.Auth(), api.AppInfo)
	r.GET("/api/get/siteinfo", api.SiteInfo)
	// 设置配置
	r.POST("/api/set/config", middleware.Auth(), api.SetConfig)
	// 清空整个表
	r.POST("/api/delete/all", middleware.Auth(), api.DeleteAll)
	// 插入数据
	r.POST("/api/add/post", middleware.Auth(), api.AddPost)
	r.POST("/api/delete/post", middleware.Auth(), api.DeletePost)

	// 初始化
	r.POST("/api/init", api.InitUser)
	// 管理员登录
	r.POST("/api/login", api.Login)
	// 检测用户是否登录
	r.GET("/api/user/is_login", middleware.Auth(), api.IsLogin)

	//获取服务端配置
	port := ":" + viper.GetString("server.port")
	// 运行服务
	r.Run(port)
}
