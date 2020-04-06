package miniprogram

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"neteaseMusicAPI/db"
	"time"
)

var appId = "wx5c201194e1d8e3a3"
var appSecret = "0eb0e5a3805ab207e956889c4a4d3e5c"

func LoginView(c *gin.Context) {
	// 获取 openid 和 session_key
	var wxLoginBaseApi = "https://api.weixin.qq.com/sns/jscode2session"
	code := c.Query("code")
	params := url.Values{}
	params.Set("appid", appId)
	params.Set("secret", appSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")
	wxLoginUrl, _ := url.Parse(wxLoginBaseApi)
	wxLoginUrl.RawQuery = params.Encode()
	wxLoginUrlPath := wxLoginUrl.String()
	resp, err := http.Get(wxLoginUrlPath)
	if err != nil {
		fmt.Println("获取登录信息失败", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	fmt.Println("body:", bodyStr)

	errcodeExist := gjson.Get(bodyStr, "errcode").Exists()
	if errcodeExist {
		c.JSON(503, gin.H{
			"code": 503,
			"msg":  gjson.Get(bodyStr, "errmsg").String(),
		})
		c.Abort()
		return
	}
	sessionKey := gjson.Get(bodyStr, "session_key").String()
	openid := gjson.Get(bodyStr, "openid").String()

	// 记录session
	sessionId := fmt.Sprintf("%x", md5.Sum([]byte(sessionKey+openid)))
	redisClient, err := db.GetRedisClient()
	if err != nil {
		fmt.Println("连接redis错误：", err)
	}

	err = redisClient.Set(sessionId, bodyStr, time.Hour*24*7).Err()
	if err != nil {
		fmt.Println("添加session失败：", err)
	}

	c.JSON(200, gin.H{
		"code":      200,
		"msg":       "登录成功",
		"sessionId": sessionId,
	})

}

func PlaylistsView(c *gin.Context) {

}

func CreatePlaylistView(c *gin.Context) {

}
func DeletePlaylistView(c *gin.Context) {

}
func GetPlaylistView(c *gin.Context) {

}

func AddSongToPlaylistView(c *gin.Context) {

}
func DelSongFromPlaylistView(c *gin.Context) {

}
