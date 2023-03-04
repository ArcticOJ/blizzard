package blizzard

import (
	"backend/blizzard/core"
	"backend/blizzard/db/seed"
	"backend/blizzard/middlewares"
	"backend/blizzard/models"
	"backend/blizzard/pb"
	"backend/blizzard/routes/auth"
	"backend/blizzard/routes/contests"
	"backend/blizzard/routes/feeds"
	"backend/blizzard/routes/problems"
	"backend/blizzard/routes/root"
	"backend/blizzard/routes/user"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oiime/logrusbun"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.arsenm.dev/drpc/muxconn"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"time"
)

var Map = map[string]models.RouteMap{
	"/problems": problems.Map,
	"/feeds":    feeds.Map,
	"/contests": contests.Map,
	"/auth":     auth.Map,
	"/user":     user.Map,
	"/":         root.Map,
}

func createHandler(server *models.Server, handler models.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := handler(&models.Context{
			Server:  server,
			Context: c,
		})
		if c.Response().Committed {
			return nil
		}
		if res == nil {
			return c.NoContent(http.StatusNoContent)
		} else {
			return c.JSONPretty(res.StatusCode(), res.Body(), "\t")
		}
	}
}

func createDb(config core.DatabaseConfig, debug bool) *bun.DB {
	psql := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithUser(config.Username),
		pgdriver.WithPassword(config.Password),
		pgdriver.WithDatabase(config.DatabaseName),
		pgdriver.WithAddr(config.Address),
		pgdriver.WithInsecure(true)))
	db := bun.NewDB(psql, pgdialect.New())
	if debug {
		log := logrus.New()
		db.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{
			Logger: log,
		}))
	}
	return db
}

func initClients(addrs map[string]string) (clients map[string]models.PolarClient) {
	clients = make(map[string]models.PolarClient)
	for name, addr := range addrs {
		// TODO: Periodically check whether the server is online and update
		clients[name] = models.PolarClient{
			DRPCPolarClient: nil,
			Address:         addr,
		}
		dial, e := net.DialTimeout("tcp", addr, time.Second*3)
		if e != nil {
			logrus.Error(e)
			continue
		}
		conn, e := muxconn.New(dial)
		if e != nil {
			logrus.Error(e)
		}
		if client, ok := clients[name]; ok {
			client.DRPCPolarClient = pb.NewDRPCPolarClient(conn)
			clients[name] = client
		}
	}
	for name := range clients {
		fmt.Println(name)
	}
	return
}

func CreateServer(config *core.Config) (server *models.Server) {
	e := echo.New()
	bootTimestamp := time.Now()
	server = &models.Server{
		Echo:          e,
		Database:      createDb(config.Database, config.Debug),
		BootTimestamp: bootTimestamp,
		Config:        config,
		Polar:         initClients(config.Judges),
	}
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code, message := http.StatusInternalServerError, "Internal Server Error"
		if er, ok := err.(*echo.HTTPError); ok {
			code = er.Code
			message = er.Message.(string)
		}
		_ = c.JSONPretty(code, models.Error{Code: code, Message: message}, "\t")
	}
	seed.Populate(server.Database)
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.RateLimit))))
	e.Pre(middleware.RemoveTrailingSlash())
	if config.EnableCORS {
		e.Use(middleware.CORS())
	}
	if config.Debug {
		e.Use(middleware.Logger())
	}
	e.Use(middlewares.Authentication(config.PrivateKey, server))
	for route, group := range Map {
		g := e.Group(route)
		for r, sub := range group {
			for _, m := range sub.Methods {
				method := m.ToString()
				handler := createHandler(server, sub.Handler)
				if route == "/" {
					e.Add(method, r, handler)
				} else if r == "/" {
					e.Add(method, route, handler)
				} else {
					g.Add(method, r, handler)
				}
			}
		}
	}
	return
}
