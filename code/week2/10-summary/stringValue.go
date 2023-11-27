package summary

import "strconv"

// StringValue 用于承载来自各部分输入的值 并提供统一的类型转换API
type StringValue struct {
	// value 承载来自各部分输入的值 以字符串表示
	value string
	// err 用于承载处理各部分输入时的错误信息
	err error
}

// AsInt64 将承载的值转换为int64类型表示
func (s StringValue) AsInt64() (value int64, err error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseInt(s.value, 10, 64)
}

// AsUint64 将承载的值转换为uint64类型表示
func (s StringValue) AsUint64() (value uint64, err error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseUint(s.value, 10, 64)
}

// AsFloat64 将承载的值转换为float64类型表示
func (s StringValue) AsFloat64() (value float64, err error) {
	if s.err != nil {
		return 0, s.err
	}

	return strconv.ParseFloat(s.value, 64)
}

//// StringValue 用于承载来自各部分输入的值 并提供统一的类型转换API
//type StringValue[T any] struct {
//	// value 承载来自各部分输入的值 以字符串表示
//	value string
//	// err 用于承载处理各部分输入时的错误信息
//	err error
//}
//
//// AsInt64 将承载的值转换为int64类型表示
//func (s StringValue[T]) AsInt64() (t T, err error) {
//	if s.err != nil {
//		return any(0), s.err
//	}
//
//	value, err := strconv.ParseInt(s.value, 10, 64)
//	if err != nil {
//		return any(0), err
//	}
//
//	return any(value), nil
//}
//
//// AsUint64 将承载的值转换为uint64类型表示
//func (s StringValue[T]) AsUint64() (t T, err error) {
//	if s.err != nil {
//		return any(0), s.err
//	}
//
//	value, err := strconv.ParseUint(s.value, 10, 64)
//	if err != nil {
//		return any(0), err
//	}
//
//	return any(value), nil
//}
//
//// AsFloat64 将承载的值转换为float64类型表示
//func (s StringValue[T]) AsFloat64() (t T, err error) {
//	if s.err != nil {
//		return any(0), s.err
//	}
//
//	value, err := strconv.ParseFloat(s.value, 64)
//	if err != nil {
//		return any(0), err
//	}
//
//	return any(value), nil
//}
