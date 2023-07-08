package extra

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/user"
	"blizzard/blizzard/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"net/http"
	"time"
)

type Context struct {
	echo.Context
}

func (ctx Context) Err(code int, message string, context ...interface{}) models.Response {
	var c interface{} = context
	if len(context) == 0 {
		c = nil
	} else if len(context) == 1 {
		c = context[0]
	}
	return &models.Error{Code: code, Message: message, Context: c}
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

func (ctx Context) Bad(message string, context ...interface{}) models.Response {
	return ctx.Err(400, message, context...)
}

func (ctx Context) Unauthorized(context ...interface{}) models.Response {
	return ctx.Err(401, "Unauthorized.", context...)
}

func (ctx Context) Forbid(message string, context ...interface{}) models.Response {
	return ctx.Err(403, message, context...)
}

func (ctx Context) NotFound(message string, context ...interface{}) models.Response {
	return ctx.Err(404, message, context...)
}

func (ctx Context) InternalServerError(message string, context ...interface{}) models.Response {
	return ctx.Err(500, message, context...)
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
	// TODO: add secure property to config
	//cookie.Secure = true
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

func (ctx Context) CommitResponse(res models.Response) error {
	return ctx.JSONPretty(res.StatusCode(), res.Body(), "\t")
}

func (ctx Context) Authenticate(uuid uuid.UUID, remember bool) models.Response {
	key := []byte(config.Config.PrivateKey)
	now := time.Now()
	lifespan := now.AddDate(0, 1, 0)
	ss := &models.Session{
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

func (ctx Context) RequireAuth() bool {
	authenticated := ctx.GetUUID() != nil
	if !authenticated {
		ctx.CommitResponse(ctx.Unauthorized())
	}
	return !authenticated
}
