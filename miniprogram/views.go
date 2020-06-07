package miniprogram

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"neteaseMusicAPI/db"
)

func endReq(c *gin.Context, code int, msg string, err interface{}) {
	fmt.Println(msg, err)

	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})
	c.Abort()
}

func LoginView(c *gin.Context) {
	// 从前端请求中获取 code 参数
	code := c.Query("code")

	// 拼接微信后台登录地址
	params := url.Values{}
	params.Set("appid", appId)
	params.Set("secret", appSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")
	wxLoginUrl, _ := url.Parse(wxLoginBaseApi)
	wxLoginUrl.RawQuery = params.Encode()
	wxLoginUrlPath := wxLoginUrl.String()

	// 获取用户 openid 和 session_key
	resp, err := http.Get(wxLoginUrlPath)
	if err != nil {
		endReq(c, 500, "请求微信后台错误", err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	bodyStr := string(body)
	if gjson.Get(bodyStr, "errcode").Exists() {
		endReq(c, 500, "微信后台返回错误", gjson.Get(bodyStr, "errmsg").String())
		return
	}
	sessionKey := gjson.Get(bodyStr, "session_key").String()
	openid := gjson.Get(bodyStr, "openid").String()

	// 生成自己 session
	sessionId := fmt.Sprintf("%x", md5.Sum([]byte(sessionKey+openid)))

	redisClient, err := db.GetRedisClient()
	if err != nil {
		endReq(c, 500, "连接redis错误", err)
		return
	}

	redisClient.Set(sessionId, bodyStr, time.Hour*24*7)

	c.JSON(200, gin.H{
		"code":      200,
		"msg":       "登录成功",
		"sessionId": sessionId,
	})

}

// get self playlist
func PlaylistsView(c *gin.Context) {
	sessionId := c.PostForm("sessionId")
	if sessionId == "" {
		endReq(c, 400, "需要登录", sessionId)
		return
	}
	redisClient, err := db.GetRedisClient()
	if err != nil {
		endReq(c, 500, "连接redis错误：", err)
		return
	}

	session := redisClient.Get(sessionId).Val()
	if session == "" {
		endReq(c, 503, "登录已过期", sessionId)
		return
	}
	openid := gjson.Get(session, "openid").String()

	playlists_str := redisClient.HGetAll(openid).Val()

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": playlists_str,
	})

}

func CreatePlaylistView(c *gin.Context) {
	//r ,err:= ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(r))
	//fmt.Println(c.Request.Header)

	sessionId := c.PostForm("sessionId")
	name := c.PostForm("name")
	desc := c.PostForm("desc")
	//fmt.Println(sessionId, name, desc)

	redisClient, err := db.GetRedisClient()
	if err != nil {
		endReq(c, 500, "连接redis错误：", err)
		return
	}

	session := redisClient.Get(sessionId).Val()
	if session == "" {
		endReq(c, 503, "登录已过期", sessionId)
		return
	}
	openid := gjson.Get(session, "openid").String()

	playlistId := name

	newPlaylist := Playlist{Name: name, Desc: desc, Songs: []map[string]string{}}
	newPlaylistByte, err := json.Marshal(newPlaylist)
	if err != nil {
		fmt.Println(err)
	}

	err = redisClient.HSet(openid, playlistId, newPlaylistByte).Err()
	if err != nil {
		fmt.Println("新建歌单失败")
		endReq(c, 500, "新建歌单失败", err)
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

	redisClient, err := db.GetRedisClient()
	if err != nil {
		endReq(c, 500, "连接redis错误：", err)
		return
	}

	session := redisClient.Get(sessionId).Val()
	if session == "" {
		endReq(c, 503, "登录已过期", sessionId)
		return
	}
	openid := gjson.Get(session, "openid").String()

	err = redisClient.HDel(openid, playlistId).Err()
	if err != nil {
		endReq(c, 500, "删除歌单失败", err)
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

	redisClient, err := db.GetRedisClient()
	if err != nil {
		endReq(c, 500, "连接redis错误：", err)
		return
	}

	session := redisClient.Get(sessionId).Val()
	if session == "" {
		endReq(c, 503, "登录已过期", sessionId)
		return
	}
	openid := gjson.Get(session, "openid").String()

	playlistByte, _ := redisClient.HGet(openid, playlistId).Bytes()

	var playlist Playlist
	err = json.Unmarshal(playlistByte, &playlist)
	if err != nil {
		endReq(c, 500, "歌单反序列化失败", err)
		return
	}

	// 添加歌曲
	playlist.Songs = append(playlist.Songs, map[string]string{"name": songName, "url": songUrl})

	newPlaylistByte, err := json.Marshal(playlist)
	resp := redisClient.HSet(openid, playlistId, string(newPlaylistByte))
	if resp.Err() != nil {
		// fmt.Println(resp.Val())
		endReq(c, 500, "保存到redis失败", resp.Err())
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

	redisClient, err := db.GetRedisClient()
	if err != nil {
		endReq(c, 500, "连接redis错误", err)
		return
	}

	session := redisClient.Get(sessionId).Val()
	if session == "" {
		endReq(c, 503, "登录已过期", sessionId)
		return
	}
	openid := gjson.Get(session, "openid").String()

	playlistByte, _ := redisClient.HGet(openid, playlistId).Bytes()

	var playlist Playlist
	err = json.Unmarshal(playlistByte, &playlist)
	if err != nil {
		endReq(c, 500, "歌单反序列化失败", err)
		return
	}

	// 删除歌曲
	for index, value := range playlist.Songs {
		if value["name"] == songName {
			playlist.Songs = append(playlist.Songs[:index], playlist.Songs[index+1:]...)
		}
	}

	newPlaylistByte, err := json.Marshal(playlist)
	redisClient.HSet(openid, playlistId, string(newPlaylistByte))

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
