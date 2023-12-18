package bindJSON

type Config struct {
	// UseNumber 反序列化JSON时是否使用Number类型
	UseNumber bool
	// DisallowUnknownFields 反序列化JSON时是否禁止未知字段
	DisallowUnknownFields bool
}
