package http

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/yudppp/throttle"
	"time"
)

type ResponseStream struct {
	encoder  *json.Encoder
	stream   *echo.Response
	throttle throttle.Throttler
}

func NewStream(response *echo.Response, interval time.Duration) *ResponseStream {
	return &ResponseStream{
		encoder:  json.NewEncoder(response),
		stream:   response,
		throttle: throttle.New(interval),
	}
}

func (rs *ResponseStream) Write(obj interface{}) error {
	e := rs.encoder.Encode(obj)
	if e != nil {
		return e
	}
	rs.throttle.Do(rs.Flush)
	return nil
}

func (rs *ResponseStream) Flush() {
	rs.stream.Flush()
}
