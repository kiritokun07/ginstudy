package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
import "github.com/thinkerou/favicon"

// 自定义Go中间件
func myHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("token", "123456")
		context.Next() //放行
		//context.Abort() //阻止
	}
}

func main() {

	//创建一个服务
	ginServer := gin.Default()
	ginServer.Use(myHandler())
	ginServer.Use(favicon.New("./favicon.ico"))

	//加载静态页面
	ginServer.LoadHTMLGlob("templates/*")

	//加载资源文件
	ginServer.Static("/static", "./static")

	//（略）连接数据库

	//访问地址
	ginServer.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"msg": "hello world",
		})
	})

	ginServer.GET("/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "这是后端传的数据",
		})
	})

	//接收前端传来的参数
	//url?userId=xxx&username=xxx
	ginServer.GET("/user/info1", myHandler(), func(context *gin.Context) {
		token := context.MustGet("token")
		log.Println("token=", token)
		userId := context.Query("userId")
		username := context.Query("username")
		context.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"username": username,
		})
	})

	// url?user/info/1/kirito
	ginServer.GET("/user/info2/:userId/:username", func(context *gin.Context) {
		userId := context.Param("userId")
		username := context.Param("username")
		context.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"username": username,
		})
	})

	//前端传json
	//curl --location --request POST 'http://localhost:8080/json' \
	//--header 'Content-Type: application/json' \
	//--data-raw '{
	//    "a": 1,
	//    "b": 2
	//}'
	ginServer.POST("/json", func(context *gin.Context) {
		data, _ := context.GetRawData()
		var m map[string]interface{}
		_ = json.Unmarshal(data, &m)
		context.JSON(http.StatusOK, m)
	})

	//路由 301
	ginServer.GET("/test", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	//路由组
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/add")
		userGroup.POST("/login")
		userGroup.DELETE("/logout")
	}

	orderGroup := ginServer.Group("/order")
	{
		orderGroup.GET("/add")
		orderGroup.DELETE("/del")
	}

	//服务器端口
	err := ginServer.Run(":8080")
	if err != nil {
		return
	}

}
