package http

type Handler func(ctx *Context) Response

type RouteMap map[string]Route

type Route struct {
	Methods []Method
	Handler
}
