package netease

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// 根据 歌曲id 从网易接口获取 歌曲url
func getSongUrlFromApi(songId string) (string, error) {
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
func getSongIdsFromApi(playlistId string) (songIds []string, err error) {

	url := baseUrl + "?type=playlist" + "&id=" + playlistId

	playlistRes, err := http.Get(url)
	if err != nil {
		return songIds, err
	}

	body, err := ioutil.ReadAll(playlistRes.Body)
	_ = playlistRes.Body.Close()

	if err != nil {
		return songIds, err
	}

	songIds = make([]string, 0)
	value := gjson.Get(string(body), "playlist.trackIds.#.id")
	for _, songId := range value.Array() {
		songIds = append(songIds, songId.String())
	}

	return songIds, nil
}
