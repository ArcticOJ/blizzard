package http

import (
	"blizzard/config"
	"blizzard/db"
	"blizzard/db/models/user"
	"crypto/md5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tmthrgd/go-hex"
	"github.com/uptrace/bun"
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

func (ctx Context) StreamResponse() *ResponseStream {
	r := ctx.Response()
	h := r.Header()
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "keep-alive")
	r.WriteHeader(http.StatusOK)
	return NewStream(r)
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

func (ctx Context) GetUser() *user.User {
	return ctx.GetDetailedUser(func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Column("id", "handle", "display_name", "email")
	})
}

func (ctx Context) GetDetailedUser(q func(query *bun.SelectQuery) *bun.SelectQuery) *user.User {
	id := ctx.GetUUID()
	if id == nil {
		return nil
	}
	var usr user.User
	query := db.Database.NewSelect().Model(&usr).Where("id = ?", id)
	if q != nil {
		query = q(query)
	}
	e := query.Scan(ctx.Request().Context())
	if e != nil {
		return nil
	}
	return &usr
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

func (ctx Context) CommitResponse(res Response) error {
	return ctx.JSON(res.StatusCode(), res.Body())
}

func (ctx Context) Authenticate(uuid uuid.UUID, remember bool) Response {
	key := []byte(config.Config.PrivateKey)
	now := time.Now()
	lifespan := now.AddDate(0, 1, 0)
	ss := &Session{
		UUID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(lifespan),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "Arctic Judge Platform",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ss)
	signedToken, e := token.SignedString(key)
	if e != nil {
		return ctx.InternalServerError("Could not create a new session.")
	}
	ctx.PutCookie("session", signedToken, lifespan, !remember)
	return ctx.Success()
}

func (ctx Context) AddAvatar(usr *user.User) {
	h := md5.Sum([]byte(usr.Email))
	usr.Avatar = hex.EncodeToString(h[:])
	usr.Email = ""
}

func (ctx Context) RequireAuth() bool {
	authenticated := ctx.GetUUID() != nil
	if !authenticated {
		ctx.CommitResponse(ctx.Unauthorized())
	}
	return !authenticated
}
