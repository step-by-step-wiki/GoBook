package v6

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
)

// SafeContext 使用装饰器模式 为 Context 添加了一个互斥锁
// 实现了 Context 的线程安全
type SafeContext struct {
	Context Context
	Lock    sync.Mutex
}

// SetCookie 设置响应头中的Set-Cookie字段 该方法是线程安全的
func (s *SafeContext) SetCookie(cookie *http.Cookie) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	http.SetCookie(s.Context.Resp, cookie)
}

// BindJSON 绑定请求体中的JSON到给定的实例(这里的实例不一定是结构体实例,还有可能是个map)上
// 该方法是线程安全的
func (s *SafeContext) BindJSON(target any) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if target == nil {
		return errors.New("web绑定错误: 给定的实例为空")
	}

	if s.Context.Req.Body == nil {
		return errors.New("web绑定错误: 请求体为空")
	}

	decoder := json.NewDecoder(s.Context.Req.Body)
	return decoder.Decode(target)
}

// FormValue 获取表单中给定键的值 该方法是线程安全的
func (s *SafeContext) FormValue(key string) (stringValue StringValue) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	err := s.Context.Req.ParseForm()
	if err != nil {
		return StringValue{err: err}
	}

	return StringValue{value: s.Context.Req.FormValue(key)}
}

// QueryValue 获取查询字符串中给定键的值 该方法是线程安全的
func (s *SafeContext) QueryValue(key string) (stringValue StringValue) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.Context.queryValues == nil {
		s.Context.queryValues = s.Context.Req.URL.Query()
	}

	if len(s.Context.queryValues) == 0 {
		return StringValue{err: errors.New("web绑定错误: 无任何查询参数")}
	}

	values, ok := s.Context.queryValues[key]
	if !ok {
		return StringValue{err: errors.New("web绑定错误: 查询参数中不存在键: " + key)}
	}

	return StringValue{value: values[0]}
}

// PathValue 获取路径参数中给定键的值 该方法是线程安全的
func (s *SafeContext) PathValue(key string) (stringValue StringValue) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.Context.PathParams == nil {
		return StringValue{err: errors.New("web绑定错误: 无任何路径参数")}
	}

	value, ok := s.Context.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web绑定错误: 路径参数中不存在键: " + key)}
	}

	return StringValue{value: value}
}

// RespJSON 以JSON格式输出相应 该方法是线程安全的
func (s *SafeContext) RespJSON(status int, obj any) (err error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	s.Context.Resp.Header().Set("Content-Type", "application/json")
	s.Context.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	// Tips: 在写入响应状态码之前设置响应头 因为一旦调用了WriteHeader方法
	// Tips: 随后对响应头的任何修改都不会生效 因为响应头已经发送给客户端了
	s.Context.Resp.WriteHeader(status)

	n, err := s.Context.Resp.Write(data)
	if n != len(data) {
		return errors.New("web绑定错误: 写入响应体不完整")
	}

	return err
}

// RespJSONOK 以JSON格式输出一个状态码为200的响应 该方法是线程安全的
func (s *SafeContext) RespJSONOK(obj any) (err error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	return s.Context.RespJSON(http.StatusOK, obj)
}
