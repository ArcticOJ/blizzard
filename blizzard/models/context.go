package models

import (
	"backend/blizzard/db/models/shared"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Context struct {
	echo.Context
	Server *Server
}

func (ctx Context) Err(code int, message string, context *interface{}) Response {
	return &Error{Code: code, Message: message, Context: context}
}

func (ctx Context) Respond(data interface{}) Response {
	return &JsonResponse{200, data}
}

func (ctx Context) Method() Method {
	return MethodFromString(ctx.Request().Method)
}

func (ctx Context) Arr(arr ...interface{}) Response {
	return ctx.Respond(arr)
}

func (ctx Context) StreamResponse() *ResponseStream {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.Response().WriteHeader(http.StatusOK)
	return New(ctx.Response())
}

func (ctx Context) Bad(message string) Response {
	return ctx.Err(400, message, nil)
}

func (ctx Context) Unauthorized() Response {
	return ctx.Err(401, "Unauthorized.", nil)
}

func (ctx Context) Forbid(message string) Response {
	return ctx.Err(403, message, nil)
}

func (ctx Context) NotFound(message string) Response {
	return ctx.Err(404, message, nil)
}

func (ctx Context) InternalServerError(message string) Response {
	return ctx.Err(500, message, nil)
}

func (ctx Context) Success() Response {
	return ctx.Respond(echo.Map{
		"success": true,
	})
}

func (ctx Context) GetUUID() *uuid.UUID {
	id := ctx.Get("user")
	if id == nil {
		return nil
	}
	if uid, ok := id.(uuid.UUID); !ok {
		return nil
	} else {
		return &uid
	}
}

func (ctx Context) GetUser(columns ...string) *shared.User {
	id := ctx.GetUUID()
	if id == nil {
		return nil
	}
	var user shared.User
	query := ctx.Server.Database.NewSelect().Model(&user).Where("id = ?", id)
	if len(columns) == 0 {
		query = query.ExcludeColumn("password", "apiKey")
	} else {
		query = query.Column(columns...)
	}
	e := query.Scan(ctx.Request().Context())
	if e != nil {
		return nil
	}
	return &user
}

func (ctx Context) PutCookie(name string, value string, exp time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = exp
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"
	ctx.SetCookie(cookie)
}

func (ctx Context) DeleteCookie(name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"
	ctx.SetCookie(cookie)
}
