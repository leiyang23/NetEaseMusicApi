package main

import (
	"io"
	"log"
	"neteaseMusicAPI/assert"
	"neteaseMusicAPI/miniprogram"
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

	groupMiniprogram := r.Group("/miniprogram")
	{
		groupMiniprogram.GET("/login", miniprogram.LoginView)
		groupMiniprogram.POST("/playlists", miniprogram.PlaylistsView)
		groupMiniprogram.POST("/playlist/create", miniprogram.CreatePlaylistView)
		groupMiniprogram.POST("/playlist/delete", miniprogram.DeletePlaylistView)
		groupMiniprogram.POST("/playlist/addSong", miniprogram.AddSongToPlaylistView)
		groupMiniprogram.POST("/playlist/delSong", miniprogram.DelSongFromPlaylistView)
	}

	r.GET("/assert", assert.AssertsView)
	r.GET("/assert/list", assert.ListView)

	r.GET("/netease/random", assert.RandomView)

	r.GET("/help", func(c *gin.Context) {
		c.HTML(200, "help.html", gin.H{
			"title": "帮助",
		})
	})

	log.Fatal(r.Run(":1627"))
}
