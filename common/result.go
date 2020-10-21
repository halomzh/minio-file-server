package common

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Result) GetMessage() string {
	return r.Message
}

func (r *Result) GetData() interface{} {
	return r.Data
}

func (r *Result) SetMessage(message string) *Result {
	r.Message = message
	return r
}

func (r *Result) SetData(data interface{}) *Result {
	r.Data = data
	return r
}

func GenSuccessResult() *Result {
	return &Result{
		Code:    0,
		Message: "success",
	}
}

func GenFailResult() *Result {
	return &Result{
		Code:    500,
		Message: "未知异常，请联系管理员",
	}
}
