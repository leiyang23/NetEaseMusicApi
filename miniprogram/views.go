package miniprogram

import (
	"crypto/md5"
	"encoding/json"
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
	sessionId := c.PostForm("sessionId")
	if sessionId == "" {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  "need session id",
		})
		c.Abort()
		return
	}
	redisClient, err := db.GetRedisClient()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "redis connection error",
		})
		c.Abort()
		return
	}

	sessionExist := redisClient.Exists(sessionId).Name()
	if sessionExist != "exists" {
		c.JSON(503, gin.H{
			"code": 503,
			"msg":  "login expired",
		})
		c.Abort()
		return
	}

	res := redisClient.Get(sessionId).Val()
	openid := gjson.Get(res, "openid").String()
	if openid == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "no record",
		})
		c.Abort()
		return
	}

	playlists_str := redisClient.HGetAll(openid).Val()
	print(playlists_str)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": playlists_str,
	})

}

type Playlist struct {
	Name  string
	Desc  string
	Songs []map[string]string
}

func CreatePlaylistView(c *gin.Context) {
	//r ,err:= ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(r))
	//fmt.Println(c.Request.Header)

	sessionId := c.PostForm("sessionId")
	name := c.PostForm("name")
	desc := c.PostForm("desc")
	fmt.Println(sessionId, name, desc)

	if sessionId == "" {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  "need session id",
		})
		c.Abort()
		return
	}
	redisClient, err := db.GetRedisClient()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "redis connection error",
		})
		c.Abort()
		return
	}

	sessionExist := redisClient.Exists(sessionId).Name()
	if sessionExist != "exists" {
		c.JSON(503, gin.H{
			"code": 503,
			"msg":  "login expired",
		})
		c.Abort()
		return
	}
	res := redisClient.Get(sessionId).Val()

	openid := gjson.Get(res, "openid").String()
	//playlistId := strconv.FormatInt(time.Now().Unix(),10)
	playlistId := name

	newPlaylist := Playlist{Name: name, Desc: desc, Songs: []map[string]string{}}
	newPlaylistByte, err := json.Marshal(newPlaylist)
	if err != nil {
		fmt.Println(err)
	}

	resp := redisClient.HSet(openid, playlistId, newPlaylistByte).Err()
	fmt.Println("新建结果：", resp)
	if resp != nil {
		fmt.Println("新建歌单失败")
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "fail to create playlist",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"code":       200,
		"msg":        "success",
		"playlistId": playlistId,
	})

}
func DeletePlaylistView(c *gin.Context) {
	sessionId := c.PostForm("sessionId")
	playlistId := c.PostForm("playlistId")

	if sessionId == "" {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  "need session id",
		})
		c.Abort()
		return
	}
	redisClient, err := db.GetRedisClient()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "redis connection error",
		})
		c.Abort()
		return
	}

	sessionExist := redisClient.Exists(sessionId).Name()
	if sessionExist != "exists" {
		c.JSON(503, gin.H{
			"code": 503,
			"msg":  "login expired",
		})
		c.Abort()
		return
	}
	res := redisClient.Get(sessionId).Val()

	openid := gjson.Get(res, "openid").String()

	err = redisClient.HDel(openid, playlistId).Err()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "fail to del playlist",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func AddSongToPlaylistView(c *gin.Context) {
	sessionId := c.PostForm("sessionId")
	playlistId := c.PostForm("playlistId")

	songName := c.PostForm("songName")
	songUrl := c.PostForm("songUrl")

	if sessionId == "" {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  "need session id",
		})
		c.Abort()
		return
	}
	redisClient, err := db.GetRedisClient()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "redis connection error",
		})
		c.Abort()
		return
	}

	sessionExist := redisClient.Exists(sessionId).Name()
	if sessionExist != "exists" {
		c.JSON(503, gin.H{
			"code": 503,
			"msg":  "login expired",
		})
		c.Abort()
		return
	}
	res := redisClient.Get(sessionId).Val()

	openid := gjson.Get(res, "openid").String()

	playlistByte, _ := redisClient.HGet(openid, playlistId).Bytes()

	var playlist Playlist
	err = json.Unmarshal(playlistByte, &playlist)
	if err != nil {
		fmt.Println(err)
	}

	// 添加歌曲
	playlist.Songs = append(playlist.Songs, map[string]string{"name": songName, "url": songUrl})

	newPlaylistByte, err := json.Marshal(playlist)
	resp := redisClient.HSet(openid, playlistId, string(newPlaylistByte))
	if resp.Err() != nil {
		fmt.Println("插入失败：", resp.Err())
		fmt.Println(resp.Val())
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "fail to add song",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func DelSongFromPlaylistView(c *gin.Context) {
	sessionId := c.PostForm("sessionId")
	playlistId := c.PostForm("playlistId")

	songName := c.PostForm("songName")

	if sessionId == "" {
		c.JSON(400, gin.H{
			"code": -1,
			"msg":  "need session id",
		})
		c.Abort()
		return
	}
	redisClient, err := db.GetRedisClient()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "redis connection error",
		})
		c.Abort()
		return
	}

	sessionExist := redisClient.Exists(sessionId).Name()
	if sessionExist != "exists" {
		c.JSON(503, gin.H{
			"code": 503,
			"msg":  "login expired",
		})
		c.Abort()
		return
	}
	res := redisClient.Get(sessionId).Val()

	openid := gjson.Get(res, "openid").String()

	playlistByte, _ := redisClient.HGet(openid, playlistId).Bytes()

	var playlist Playlist
	err = json.Unmarshal(playlistByte, &playlist)
	if err != nil {
		fmt.Println(err)
	}

	// 删除歌曲
	for index, value := range playlist.Songs {
		if value["name"] == songName {
			playlist.Songs = append(playlist.Songs[:index], playlist.Songs[index+1:]...)
		}
	}

	newPlaylistByte, err := json.Marshal(playlist)
	if !redisClient.HSet(openid, playlistId, string(newPlaylistByte)).Val() {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "fail to del song",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
