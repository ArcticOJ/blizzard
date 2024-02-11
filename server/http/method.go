package http

import "net/http"

type Method = string

const (
	GET    Method = http.MethodGet
	POST   Method = http.MethodPost
	PATCH  Method = http.MethodPatch
	DELETE Method = http.MethodDelete
	PUT    Method = http.MethodPut
)

func MethodFromString(method string) Method {
	switch method {
	case http.MethodPost:
		return POST
	case http.MethodPatch:
		return PATCH
	case http.MethodDelete:
		return DELETE
	case http.MethodPut:
		return PUT
	}
	return GET
}
