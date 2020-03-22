package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
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

	r.GET("/assert", AssertsView)
	r.GET("/random", RandomView)
	r.GET("/help", func(c *gin.Context) {
		c.HTML(200, "help.html", gin.H{
			"title": "帮助",
		})
	})
	go TickClearPlaylistCache()
	go TickClearSongUrlCache()

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

// 本网站的资源
type AssertParam struct {
	Category string `form:"category" binding:"required"`
	Tag      string `form:"tag" binding:"required"`
}

func AssertsView(c *gin.Context) {
	var param AssertParam
	var basePath = "/home/asserts"

	if err := c.ShouldBind(&param); err != nil {
		c.String(400, "error: %s", err)
		c.Abort()
		return
	}
	aCategory := c.Query("category")
	aTag := c.Query("tag")

	assertPath := basePath + "/" + aCategory + "/" + aTag
	rd, err := ioutil.ReadDir(assertPath)
	if err != nil {
		fmt.Println("资源路径不存在")
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "资源路径不存在",
		})
	}

	var data []string
	urlBasePath := "https://asserts.freaks.group" + "/" + aCategory + "/" + aTag
	for _, fi := range rd {
		fmt.Println(fi.Name())
		data = append(data, fi.Name())
	}
	c.JSON(200, gin.H{
		"code":        200,
		"msg":         "success",
		"data":        data,
		"urlBasePath": urlBasePath,
	})
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
