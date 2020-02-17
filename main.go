package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()

	// 设置生产模式，默认 debug 模式
	gin.SetMode(gin.ReleaseMode)

	// 记录到文件。
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// 开始
	r := gin.New()
	r.Use(gin.Logger())

	r.LoadHTMLFiles("templates/help.html")

	r.GET("/random", RandomView)
	r.GET("/help", func(c *gin.Context) {
		c.HTML(200, "help.html", gin.H{
			"title": "帮助",
		})
	})
	go TickClearCache()

	log.Fatal(r.Run(":1627"))
}

func RandomView(c *gin.Context) {

	PlaylistId := c.DefaultQuery("playlist_id", "2467106683")
	res, err := Random(PlaylistId)

	statusCode := 200
	if err != nil {
		statusCode = 404
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.String(statusCode, res)

}

// 参数格式
type ResourceParam struct {
	Type string `form:"type" binding:"required"`
	Id   string `form:"id" binding:"required"`
}

func ResourceView(c *gin.Context) {
	// 验证参数
	var param ResourceParam
	if err := c.ShouldBind(&param); err != nil {
		c.String(400, "error: %s", err)
		c.Abort()
		return
	}

	rType := c.DefaultQuery("type", "1")
	rId := c.Query("id")
	fmt.Println(rType, rId)

	c.String(200, "nothing")
}
