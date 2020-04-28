package models

type Response struct {
	IsSuccess bool        `json:"isSuccess"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

func (r *Response) Success(data interface{}) {
	r.IsSuccess = true
	r.Message = "操作成功"
	r.Data = data
}

func (r *Response) Failed(msg string) {
	r.IsSuccess = false
	r.Message = msg
}

type APIResponse struct {
	Result int         `json:"result"`  // 状态码 0 成功， 1 请求错误， 2 条件错误， 3 内部错误
	Msg    string      `json:"message"` // 返回的消息
	Data   interface{} `json:"data"`
}

func (ar *APIResponse) Success(data interface{}) {
	ar.Data = data
	ar.Msg = "操作成功"
}

func (ar *APIResponse) BadRequest(msg string) {
	ar.Result = 1
	ar.Msg = msg
}

func (ar *APIResponse) PreCheckError(msg string) {
	ar.Result = 2
	ar.Msg = msg
}

func (ar *APIResponse) PreCheckErrorWithResult(result int, msg string) {
	ar.Result = result
	ar.Msg = msg
}

func (ar *APIResponse) InternalError(msg string) {
	ar.Result = 3
	ar.Msg = msg
}
