package miniprogram

var appId = "wx5c201194e1d8e3a3"
var appSecret = "0eb0e5a3805ab207e956889c4a4d3e5c"

var wxLoginBaseApi = "https://api.weixin.qq.com/sns/jscode2session"

type Playlist struct {
	Name  string
	Desc  string
	Songs []map[string]string
}
