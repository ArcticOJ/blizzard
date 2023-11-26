package http

import "net/http"

type Method = string

const (
	Get    Method = http.MethodGet
	Post          = http.MethodPost
	Patch         = http.MethodPatch
	Delete        = http.MethodDelete
	Put           = http.MethodPut
)

func MethodFromString(method string) Method {
	switch method {
	case http.MethodPost:
		return Post
	case http.MethodPatch:
		return Patch
	case http.MethodDelete:
		return Delete
	case http.MethodPut:
		return Put
	}
	return Get
}
