package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	r.GET("/cloudMusic", CloudMusicView)
	r.GET("/cloudMusic/resource", ResourceView)
	r.GET("/cloudMusic/search", SearchView)
	r.GET("/help", func(c *gin.Context) {
		c.HTML(200, "help.html", gin.H{
			"title": "帮助",
		})
	})

	log.Fatal(r.Run(":1627"))
}

func CloudMusicView(c *gin.Context) {
	url := c.Request.URL.RawQuery
	res, err := CloudMusic(url)

	statusCode := 200
	if err != nil {
		statusCode = 404
	}

	c.String(statusCode, res)

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
	res, err := GetResource(rType, rId)

	statusCode := 200
	if err != nil {
		statusCode = 404
	}

	c.String(statusCode, res)
}

func SearchView(c *gin.Context) {
	// 验证参数
	var param SearchParam
	if err := c.ShouldBind(&param); err != nil {
		c.String(400, "error: %s", err)
		c.Abort()
		return
	}

	sType := c.DefaultQuery("search_type", "1")
	sKeyword := c.Query("keyword")

	res, err := Search(sType, sKeyword)

	statusCode := 200
	if err != nil {
		statusCode = 404
	}

	c.String(statusCode, res)
}
