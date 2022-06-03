package web

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type ResponseCode struct {
	Code int   `json:"code"`
	Err  error `json:"err"`
}

func NewResponse(data interface{}) Response {
	return Response{data, ""}
}

func DecodeError(err string) Response {
	return Response{nil, err}
}

func NewCodeResponse(code int, err error) ResponseCode {
	return ResponseCode{
		code, err,
	}
}
