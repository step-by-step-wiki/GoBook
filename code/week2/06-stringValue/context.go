package stringValue

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Context HandleFunc的上下文
type Context struct {
	// Req 请求
	Req *http.Request
	// Resp 响应
	Resp http.ResponseWriter
	// PathParams 路径参数名值对
	PathParams map[string]string
	// QueryValues 查询参数名值对
	queryValues url.Values
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

// FormValue 获取表单中给定键的值
func (c *Context) FormValue(key string) (stringValue StringValue) {
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{err: err}
	}

	return StringValue{value: c.Req.FormValue(key)}
}

// QueryValue 获取查询字符串中给定键的值
func (c *Context) QueryValue(key string) (stringValue StringValue) {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	if len(c.queryValues) == 0 {
		return StringValue{err: errors.New("web绑定错误: 无任何查询参数")}
	}

	values, ok := c.queryValues[key]
	if !ok {
		return StringValue{err: errors.New("web绑定错误: 查询参数中不存在键: " + key)}
	}

	return StringValue{value: values[0]}
}

// PathValue 获取路径参数中给定键的值
func (c *Context) PathValue(key string) (stringValue StringValue) {
	if c.PathParams == nil {
		return StringValue{err: errors.New("web绑定错误: 无任何路径参数")}
	}

	value, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web绑定错误: 路径参数中不存在键: " + key)}
	}

	return StringValue{value: value}
}
