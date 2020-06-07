package local

import (
	"fmt"
	"io/ioutil"
	"runtime"
)

// 本网站的资源
var localBase = "/home/assert"
var urlBase = "https://assert.freaks.group"

func init() {
	sysType := runtime.GOOS
	if sysType == "windows" {
		// 本地调试地址，需要配合 Nginx
		localBase = "D:/home-server/assert"
		urlBase = "http://127.0.0.1"
	}
}

// 所有分类
func List() (data map[string][]string, err error) {
	rd, err := ioutil.ReadDir(localBase)

	data = make(map[string][]string)

	for _, category := range rd {
		if category.IsDir() {
			categoryName := category.Name()
			sonDir := localBase + "/" + categoryName
			rd2, _ := ioutil.ReadDir(sonDir)
			for _, tag := range rd2 {
				if tag.IsDir() {
					data[categoryName] = append(data[categoryName], tag.Name())
				}

			}
		}
	}
	return data, err
}

// 返回二级分类下的所有文件
func Detail(category, tag string) (data []string, urlBasePath string, err error) {

	urlBasePath = urlBase + "/" + category + "/" + tag

	assertPath := localBase + "/" + category + "/" + tag
	rd, err := ioutil.ReadDir(assertPath)

	for _, fi := range rd {
		data = append(data, fi.Name())
	}

	return data, urlBasePath, err
}

// 测试
func main() {
	data, err := List()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)

	data2, urlbase, err := Detail("audio", "sunyanzi")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(urlbase)
	fmt.Println(data2)
}
