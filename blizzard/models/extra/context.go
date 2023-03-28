package extra

import (
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/shared"
	"blizzard/blizzard/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Context struct {
	echo.Context
}

func (ctx Context) Err(code int, message string, context *interface{}) models.Response {
	return &models.Error{Code: code, Message: message, Context: context}
}

func (ctx Context) Respond(data interface{}) models.Response {
	return &models.JsonResponse{Code: 200, Content: data}
}

func (ctx Context) Method() models.Method {
	return models.MethodFromString(ctx.Request().Method)
}

func (ctx Context) Arr(arr ...interface{}) models.Response {
	return ctx.Respond(arr)
}

func (ctx Context) StreamResponse() *models.ResponseStream {
	r := ctx.Response()
	h := r.Header()
	h.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "keep-alive")
	r.WriteHeader(http.StatusOK)
	return models.New(r)
}

func (ctx Context) Bad(message string) models.Response {
	return ctx.Err(400, message, nil)
}

func (ctx Context) Unauthorized() models.Response {
	return ctx.Err(401, "Unauthorized.", nil)
}

func (ctx Context) Forbid(message string) models.Response {
	return ctx.Err(403, message, nil)
}

func (ctx Context) NotFound(message string) models.Response {
	return ctx.Err(404, message, nil)
}

func (ctx Context) InternalServerError(message string) models.Response {
	return ctx.Err(500, message, nil)
}

func (ctx Context) Success() models.Response {
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
	query := db.Database.NewSelect().Model(&user).Where("id = ?", id)
	if len(columns) == 0 {
		query = query.ExcludeColumn("password", "api_key")
	} else {
		query = query.Column(columns...)
	}
	e := query.Scan(ctx.Request().Context())
	if e != nil {
		return nil
	}
	return &user
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
	// TODO: add secure property to config
	cookie.Secure = true
	ctx.SetCookie(cookie)
}

func (ctx Context) DeleteCookie(name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	cookie.MaxAge = -1
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"
	ctx.SetCookie(cookie)
}
