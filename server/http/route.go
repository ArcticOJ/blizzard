package http

type (
	Handler   func(ctx *Context) Response
	RouteFlag = uint8
	Route     struct {
		Path    string
		Method  Method
		Flags   RouteFlag
		Handler Handler
	}
)

const (
	// RouteAuth Protect this auth with authentication
	RouteAuth RouteFlag = 1 << iota
)

func (r *Route) HasFlag(flag RouteFlag) bool {
	return r.Flags&flag == flag
}
