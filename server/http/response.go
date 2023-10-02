package http

type Response interface {
	StatusCode() int
	Body() interface{}
}

type JsonResponse struct {
	Code    int         `json:"-"`
	Content interface{} `json:"-"`
}

func (r JsonResponse) StatusCode() int {
	return r.Code
}

func (r JsonResponse) Body() interface{} {
	return r.Content
}
