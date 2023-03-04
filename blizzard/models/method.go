package models

import "net/http"

type Method int

const (
	Get Method = iota
	Post
	Patch
	Delete
	Put
)

func (method Method) ToString() string {
	switch method {
	case Post:
		return http.MethodPost
	case Patch:
		return http.MethodPatch
	case Delete:
		return http.MethodDelete
	case Put:
		return http.MethodPut
	}
	return http.MethodGet
}

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
