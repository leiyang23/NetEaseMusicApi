package assert

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// 随机返回一首网易歌单内的歌曲地址
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
var localBase = "/home/assert"
var urlBase = "https://assert.freaks.group"

//var localBase = "D:/home-server/assert"
//var urlBase = "http://127.0.0.1"

func ListView(c *gin.Context) {
	rd, err := ioutil.ReadDir(localBase)
	if err != nil {
		fmt.Println("资源路径不存在")
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "资源路径不存在",
		})
		c.Abort()
		return
	}

	data := make(map[string][]string)
	for _, category := range rd {
		if category.IsDir() {
			categoryName := category.Name()
			sonDir := localBase + "/" + categoryName
			rd2, _ := ioutil.ReadDir(sonDir)
			for _, tag := range rd2 {
				data[categoryName] = append(data[categoryName], tag.Name())
			}
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

type Param struct {
	Category string `form:"category" binding:"required"`
	Tag      string `form:"tag" binding:"required"`
}

func AssertsView(c *gin.Context) {
	var param Param

	if err := c.ShouldBind(&param); err != nil {
		c.String(400, "error: %s", err)
		c.Abort()
		return
	}
	aCategory := c.Query("category")
	aTag := c.Query("tag")

	assertPath := localBase + "/" + aCategory + "/" + aTag
	rd, err := ioutil.ReadDir(assertPath)
	if err != nil {
		fmt.Println("资源路径不存在")
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "资源路径不存在",
		})
		c.Abort()
		return
	}

	var data []string
	urlBasePath := urlBase + "/" + aCategory + "/" + aTag
	for _, fi := range rd {
		data = append(data, fi.Name())
	}
	c.JSON(200, gin.H{
		"code":        200,
		"msg":         "success",
		"data":        data,
		"urlBasePath": urlBasePath,
	})
}
