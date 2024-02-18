package http

import (
	"github.com/ArcticOJ/blizzard/v0/server/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Context struct {
	echo.Context
}

func (ctx Context) Err(code int, message string, context ...interface{}) Response {
	var c interface{} = context
	if len(context) == 0 {
		c = nil
	} else if len(context) == 1 {
		c = context[0]
	}
	return &Error{Code: code, Message: message, Context: c}
}

func (ctx Context) Respond(data interface{}) Response {
	return &JsonResponse{Code: 200, Content: data}
}

func (ctx Context) Method() Method {
	return MethodFromString(ctx.Request().Method)
}

func (ctx Context) Arr(arr ...interface{}) Response {
	return ctx.Respond(arr)
}

func (ctx Context) StreamResponse(flushInterval time.Duration) *ResponseStream {
	r := ctx.Response()
	h := r.Header()
	h.Set("X-Streamed", "true")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "keep-alive")
	r.WriteHeader(http.StatusOK)
	return NewStream(r, flushInterval)
}

func (ctx Context) Bad(message string, context ...interface{}) Response {
	return ctx.Err(400, message, context...)
}

func (ctx Context) Unauthorized(context ...interface{}) Response {
	return ctx.Err(401, "Unauthorized.", context...)
}

func (ctx Context) Forbid(message string, context ...interface{}) Response {
	return ctx.Err(403, message, context...)
}

func (ctx Context) NotFound(message string, context ...interface{}) Response {
	return ctx.Err(404, message, context...)
}

func (ctx Context) InternalServerError(message string, context ...interface{}) Response {
	return ctx.Err(500, message, context...)
}

func (ctx Context) Success() Response {
	return ctx.Respond(echo.Map{
		"success": true,
	})
}

func (ctx Context) GetUUID() uuid.UUID {
	id := ctx.Get("id")
	if id == nil {
		return uuid.Nil
	}
	if uid, ok := id.(uuid.UUID); !ok {
		return uuid.Nil
	} else {
		return uid
	}
}

func (ctx Context) PutCookie(name string, value string, exp time.Time, sessionOnly bool) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	if !sessionOnly {
		cookie.Expires = exp
	}
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"
	cookie.HttpOnly = true
	//cookie.Secure = true
	ctx.SetCookie(cookie)
}

func (ctx Context) DeleteCookie(name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	cookie.Path = "/"
	ctx.SetCookie(cookie)
}

func (ctx Context) CommitResponse(res Response) error {
	return ctx.JSON(res.StatusCode(), res.Body())
}

func (ctx Context) Authenticate(uid uuid.UUID, remember bool) Response {
	lifespan := time.Hour * 24
	if remember {
		lifespan *= 30
	}
	if k, validUntil := session.Encrypt(lifespan, uid); k != "" {
		ctx.PutCookie("session", k, validUntil, !remember)
		return ctx.Success()
	}
	return ctx.InternalServerError("Could not create a new session.")
}

func (ctx Context) RequireAuth() bool {
	authenticated := ctx.GetUUID() != uuid.Nil
	if !authenticated {
		ctx.CommitResponse(ctx.Unauthenticated())
	}
	return !authenticated
}

func (ctx Context) Unauthenticated() Response {
	return ctx.Err(403, "Unauthenticated.")
}
