package main

import (
	"io/ioutil"
	"net/http"
)

// 接口文档：https://zhuanlan.zhihu.com/p/30246788
var baseUrl string = "https://api.imjad.cn/cloudmusic/"

// 参数格式
type ResourceParam struct {
	Type string `form:"type" binding:"required"`
	Id   string `form:"id" binding:"required"`
}

type SearchParam struct {
	Search_type string `form:"search_type" binding:"required"`
	Keyword     string `form:"keyword" binding:"required"`
}

// 请求网易后台数据

func CloudMusic(rawQuery string) (res string, err error) {
	url := baseUrl + "?" + rawQuery

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}

func Search(searchType string, keyWord string) (res string, err error) {
	url := baseUrl + "?type=search" + "&search_type=" + searchType + "&s=" + keyWord

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func GetResource(resourceType, resourceId string) (res string, err error) {
	url := baseUrl + "?type=" + resourceType + "&id=" + resourceId

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
