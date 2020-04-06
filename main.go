package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"neteaseMusicAPI/assert"
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

	r.GET("/assert", assert.AssertsView)
	r.GET("/assert/list", assert.ListView)
	r.GET("/random", assert.RandomView)
	r.GET("/help", func(c *gin.Context) {
		c.HTML(200, "help.html", gin.H{
			"title": "帮助",
		})
	})

	// 缓存任务
	go assert.GoTickClearPlaylistCache()
	go assert.GoTickClearSongUrlCache()

	log.Fatal(r.Run(":1627"))
}
