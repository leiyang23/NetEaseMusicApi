package netease

import (
	"errors"
	"fmt"
	"math/rand"
	"neteaseMusicAPI/db"
	"sync"
	"time"
)

// 接口文档：https://zhuanlan.zhihu.com/p/30246788

var (
	baseUrl           = "https://api.imjad.cn/cloudmusic/"
	CachePlaylist     = make(map[string][]string) // 歌单缓存
	CacheSongUrl      = make(map[string]string)   // 歌曲地址缓存
	defaultPlaylistId string                      // 默认歌单ID
	mutex             sync.RWMutex
)

func init() {
	// 周期任务
	go goTickClearPlaylistCache()
	go goTickClearSongUrlCache()
	go goTickChangeDefaultPlaylistId()
}

// 定时清除缓存 歌单缓存
func goTickClearPlaylistCache() {
	c := time.Tick(6 * time.Hour)
	for {
		<-c
		fmt.Println("清除歌单缓存")
		mutex.Lock()
		for i := range CachePlaylist {
			delete(CachePlaylist, i)
		}
		mutex.Unlock()
	}
}

// 定时清除缓存 歌曲缓存
func goTickClearSongUrlCache() {
	c := time.Tick(12 * time.Hour)
	for {
		<-c
		fmt.Println("清除歌曲地址缓存")
		mutex.Lock()
		for i := range CacheSongUrl {
			delete(CacheSongUrl, i)
		}
		mutex.Unlock()
	}
}

// 定时拉取默认歌单
func goTickChangeDefaultPlaylistId() {
	redisClient, _ := db.GetRedisClient()
	defaultPlaylistId = redisClient.Get("default_playlist_id").Val()
	if defaultPlaylistId == "" {
		defaultPlaylistId = "5052261708"
	}
	redisClient.Close()
	c := time.Tick(1 * time.Hour)
	for {
		<-c
		fmt.Println("更新默认歌单Id")
		redisClient, _ := db.GetRedisClient()
		defaultPlaylistId = redisClient.Get("default_playlist_id").Val()
		mutex.Lock()
		if defaultPlaylistId == "" {
			defaultPlaylistId = "5052261708"
		}
		mutex.Unlock()
		redisClient.Close()
	}
}

// 根据歌单id 从缓存中随机取出一个 歌曲id，若不存在，就请求网易接口
func randomSongIdInCache(playlistId string) (string, error) {
	var ok bool
	var playlist []string

	playlist, ok = CachePlaylist[playlistId]
	if !ok {
		songIds, err := getSongIdsFromApi(playlistId)
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
func getUrlBySongId(songId string) (string, error) {
	songUrl, ok := CacheSongUrl[songId]
	if !ok {
		newSongUrl, err := getSongUrlFromApi(songId)
		if newSongUrl == "" {
			// 对于已下架或收费的歌曲，无法获取地址，从缓存中删除songid,防止下次取到
			delete(CacheSongUrl, songId)
			delete(CachePlaylist, songId)
			return "", nil
		}
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		CacheSongUrl[songId] = newSongUrl
		return newSongUrl, nil
	}
	return songUrl, nil
}

// 随机返回一首网易歌单内的歌曲地址
func Random(playlistId string) (songUrl string, err error) {
	if playlistId == "" {
		playlistId = defaultPlaylistId
	}
	songId, err := randomSongIdInCache(playlistId)
	if err != nil {
		return "", err
	}

	songUrl, err = getUrlBySongId(songId)

	return songUrl, err
}
