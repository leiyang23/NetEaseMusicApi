package assert

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// 接口文档：https://zhuanlan.zhihu.com/p/30246788
var baseUrl string = "https://api.imjad.cn/cloudmusic/"

// 歌单缓存
var CachePlaylist = make(map[string][]string)

// 歌曲地址缓存
var CacheSongUrl = make(map[string]string)

// 定时清除缓存任务
func GoTickClearPlaylistCache() {
	// 定时清除缓存 歌单缓存
	c := time.Tick(6 * time.Hour)
	for {
		<-c
		fmt.Println("清除歌单缓存")
		for i := range CachePlaylist {
			delete(CachePlaylist, i)
		}
	}
}
func GoTickClearSongUrlCache() {
	// 定时清除缓存 歌曲缓存
	c := time.Tick(12 * time.Hour)
	for {
		<-c
		fmt.Println("清除歌曲地址缓存")
		for i := range CacheSongUrl {
			delete(CacheSongUrl, i)
		}
	}
}

// 根据 歌曲id 从网易接口获取 歌曲url
func _getSongUrlByReq(songId string) (string, error) {
	url := baseUrl + "?type=song" + "&id=" + songId
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	songUrlResult := gjson.Get(string(body2), "data.0.url")
	songUrl := strings.ReplaceAll(songUrlResult.String(), "\\", "")
	return songUrl, nil
}

// 根据歌单id 从网易接口获取 歌单内容
func _updatePlaylistByReq(playlistId string) ([]string, error) {

	fmt.Println("更新或者初始化歌单数据：", playlistId)
	url := baseUrl + "?type=playlist" + "&id=" + playlistId
	playlistRes, err := http.Get(url)
	defer playlistRes.Body.Close()
	if err != nil {
		return []string{}, err
	}

	body, err := ioutil.ReadAll(playlistRes.Body)
	if err != nil {
		return []string{}, err
	}

	songIds := make([]string, 10)
	value := gjson.Get(string(body), "playlist.trackIds.#.id")
	for _, songId := range value.Array() {
		songIds = append(songIds, songId.String())
	}

	return songIds, nil
}

// 根据歌单id 从缓存中随机取出一个 歌曲id，若不存在，就请求网易接口
func getOneSongIdByCache(playlistId string) (string, error) {
	var ok bool
	var playlist []string
	playlist, ok = CachePlaylist[playlistId]
	if !ok {
		songIds, err := _updatePlaylistByReq(playlistId)
		if err != nil {
			return "", err
		}
		if len(songIds) <= 0 {
			return "", errors.New("歌单内无歌曲")
		}

		CachePlaylist[playlistId] = songIds
		playlist = songIds
	}

	// 随机从歌单中选取一个 歌曲id
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(playlist))

	songId := playlist[index]
	return songId, nil
}

// 根据歌曲id 从缓存中查找 歌曲url，若不存在，就请求网易接口
func getSongUrlById(songId string) (string, error) {
	songUrl, ok := CacheSongUrl[songId]
	if !ok {
		newSongUrl, err := _getSongUrlByReq(songId)
		if err != nil {
			return "", err
		}
		CacheSongUrl[songId] = newSongUrl
		return newSongUrl, nil
	}
	return songUrl, nil
}

// 随机返回一首网易歌单内的歌曲地址
func Random(playlistId string) (res string, err error) {
	songId, err := getOneSongIdByCache(playlistId)
	if err != nil {
		return "", err
	}
	fmt.Println("获取歌曲id:", songId)

	songUrl, err := getSongUrlById(songId)

	return songUrl, err
}
