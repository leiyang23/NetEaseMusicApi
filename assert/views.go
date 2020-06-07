package assert

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"neteaseMusicAPI/assert/local"
	"neteaseMusicAPI/assert/netease"
)

func endReq(c *gin.Context, code int, msg string, err interface{}) {
	fmt.Println(msg, err)

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})
	c.Abort()
}

// 随机返回一首网易歌单内的歌曲地址
func RandomView(c *gin.Context) {
	PlaylistId := c.DefaultQuery("playlist_id", "")
	res, err := netease.Random(PlaylistId)

	statusCode := 200
	if err != nil {
		statusCode = 404
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.String(statusCode, res)
}

func ListView(c *gin.Context) {
	data, err := local.List()
	if err != nil {
		endReq(c, 500, "获取列表失败：", err)
		return
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

	data, urlBasePath, err := local.Detail(aCategory, aTag)
	if err != nil {
		endReq(c, 500, "获取资源失败：", err)
		return
	}

	c.JSON(200, gin.H{
		"code":        200,
		"msg":         "success",
		"data":        data,
		"urlBasePath": urlBasePath,
	})
}
