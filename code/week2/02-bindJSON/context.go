package bindJSON

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Context HandleFunc的上下文
type Context struct {
	// Req 请求
	Req *http.Request
	// Resp 响应
	Resp http.ResponseWriter
	// PathParams 路径参数名值对
	PathParams map[string]string
}

// BindJSON 绑定请求体中的JSON到给定的实例(这里的实例不一定是结构体实例,还有可能是个map)上
func (c *Context) BindJSON(target any) error {
	if target == nil {
		return errors.New("web绑定错误: 给定的实例为空")
	}

	if c.Req.Body == nil {
		return errors.New("web绑定错误: 请求体为空")
	}

	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(target)
}

// BindJSONOpt 绑定请求体中的JSON到给定的实例(这里的实例不一定是结构体实例,还有可能是个map)上
// 同时支持指定是否使用Number类型,以及是否禁止未知字段
func (c *Context) BindJSONOpt(target any, useNumber bool, disallowUnknownFields bool) error {
	if target == nil {
		return errors.New("web绑定错误: 给定的实例为空")
	}

	if c.Req.Body == nil {
		return errors.New("web绑定错误: 请求体为空")
	}

	decoder := json.NewDecoder(c.Req.Body)

	if useNumber {
		decoder.UseNumber()
	}

	if disallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	return decoder.Decode(target)
}
