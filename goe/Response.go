package goe

var (
	OK  = NewResponse(0, "success")
	Err = NewResponse(1, "error")
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(code int, message string) *Response {
	return &Response{Code: code, Message: message, Data: nil}
}

func (this *Response) WithData(data interface{}) *Response {
	this.Data = data
	return this
}

func (this *Response) WithMessage(message string) *Response {
	this.Message = message
	return this
}

func (this *Response) WithCode(code int) *Response {
	this.Code = code
	return this
}
