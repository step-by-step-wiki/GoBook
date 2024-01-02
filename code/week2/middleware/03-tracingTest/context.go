package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Context HandleFunc的上下文
type Context struct {
	Req            *http.Request       // Req 请求
	Resp           http.ResponseWriter // Resp 响应
	PathParams     map[string]string   // PathParams 路径参数名值对
	queryValues    url.Values          // queryValues 查询参数名值对
	cookieSameSite http.SameSite       // cookieSameSite cookie的SameSite属性 即同源策略
	MatchRoute     string              // MatchRoute 命中的路由
	RespData       []byte              // RespData 响应数据 主要是给中间件使用
	RespStatusCode int                 // RespStatusCode 响应状态码 主要是给中间件使用
}

// SetCookie 设置响应头中的Set-Cookie字段
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
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

// RespJSON 以JSON格式输出相应
func (c *Context) RespJSON(status int, obj any) (err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))

	// 将响应码和响应数据写入到对应字段上 而不再是直接写入到响应中
	c.RespStatusCode = status
	c.RespData = data

	return err
}

// RespJSONOK 以JSON格式输出一个状态码为200的响应
func (c *Context) RespJSONOK(obj any) (err error) {

	return c.RespJSON(http.StatusOK, obj)
}
