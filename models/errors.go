package models

type Err struct {
	// 错误码, 0 : 无错误, 1 : 程序错误, 2 : 自定义错误
	Code int
	// 信息
	Msg string
	// 错误
	Data interface{}
}

func (er *Err) UseError(code int, msg string, v interface{}) {
	er.Data = v
	er.Code = code
	er.Msg = msg
}

func (er *Err) Error() string {
	return er.Msg
}

func NewErr(msg string) error {
	return &Err{2, msg, nil}
}

func NewErrWithCode(i int, msg string) error {
	return &Err{i, msg, nil}
}
