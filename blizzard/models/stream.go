package models

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
)

type ResponseStream struct {
	Encoder *json.Encoder
	Stream  *echo.Response
}

func New(response *echo.Response) *ResponseStream {
	return &ResponseStream{
		Encoder: json.NewEncoder(response),
		Stream:  response,
	}
}

func (stream *ResponseStream) Write(obj interface{}) error {
	e := stream.Encoder.Encode(obj)
	if e != nil {
		return e
	}
	stream.Flush()
	return nil
}

func (stream *ResponseStream) Flush() {
	stream.Stream.Flush()
}
