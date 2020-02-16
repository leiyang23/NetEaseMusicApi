package main

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

var CachePlaylist = make(map[string][]string)

func TickClearCache() {
	// 定时清除缓存
	c := time.Tick(6 * time.Hour)
	for {
		<-c
		fmt.Println("清除缓存")
		for i := range CachePlaylist {
			delete(CachePlaylist, i)
		}
	}
}

func GetOneSongUrl(playlistId string) (string, error) {
	// 随机从缓存集合中取出一个 song id
	var ok bool
	var playlist []string
	playlist, ok = CachePlaylist[playlistId]
	if !ok {
		songIds, err := updatePlaylist(playlistId)
		if err != nil {
			return "", err
		}
		if len(songIds) <= 0 {
			return "", errors.New("歌单内无歌曲")
		}

		CachePlaylist[playlistId] = songIds
		playlist = songIds
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(playlist) - 1)

	songId := playlist[index]
	return songId, nil

}

func updatePlaylist(playlistId string) ([]string, error) {
	// 根据歌单id, 获取所有歌曲id
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

// 请求网易后台数据
func Random(playlistId string) (res string, err error) {

	songId, err := GetOneSongUrl(playlistId)
	if err != nil {
		return "", err
	}
	fmt.Println("获取歌曲id:", songId)

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

	return songUrl, err
}
