package models

type Error struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Context *interface{} `json:"context,omitempty"`
}

func (e Error) StatusCode() int {
	return e.Code
}

func (e Error) Body() interface{} {
	return e
}
