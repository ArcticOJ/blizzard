package extra

import "blizzard/blizzard/models"

type Handler func(ctx *Context) models.Response

type RouteMap map[string]Route

type Route struct {
	Methods []models.Method
	Handler
}
