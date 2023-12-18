package formValue

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

// FormValue1 获取表单中给定键的值
func (c *Context) FormValue1(key string) (value string, err error) {
	err = c.Req.ParseForm()
	if err != nil {
		return "", errors.New("web绑定错误: 解析表单失败: " + err.Error())
	}

	values, ok := c.Req.Form[key]
	if !ok {
		return "", errors.New("web绑定错误: 表单中不存在键: " + key)
	}

	// Tips: 这里只返回第一个值,这样的设计是参照了net/http包中的FormValue()方法
	return values[0], nil
}

// FormValue2 获取表单中给定键的值 不建议按此方式返回 因为大部分场景下表单中的键都是唯一的
func (c *Context) FormValue2(key string) (values []string, err error) {
	err = c.Req.ParseForm()
	if err != nil {
		return nil, errors.New("web绑定错误: 解析表单失败: " + err.Error())
	}

	values, ok := c.Req.Form[key]
	if !ok {
		return nil, errors.New("web绑定错误: 表单中不存在键: " + key)
	}

	return values, nil
}

// FormValue3 获取表单中给定键的值 推荐使用这种实现 因为这种实现的语义和原生API语义相同
func (c *Context) FormValue3(key string) (value string, err error) {
	err = c.Req.ParseForm()
	if err != nil {
		return "", errors.New("web绑定错误: 解析表单失败: " + err.Error())
	}

	return c.Req.FormValue(key), nil
}

// FormValueAsInt64 获取表单中给定键的值 并将该值转换为int64类型返回
func (c *Context) FormValueAsInt64(key string) (int64Value int64, err error) {
	err = c.Req.ParseForm()
	if err != nil {
		return 0, errors.New("web绑定错误: 解析表单失败: " + err.Error())
	}

	value := c.Req.FormValue(key)
	return strconv.ParseInt(value, 10, 64)
}
