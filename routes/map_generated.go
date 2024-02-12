// Code generated by /tmp/arctic/gen_routes ./blizzard/routes github.com/ArcticOJ/blizzard/v0 ./blizzard/routes/map_generated.go; DO NOT EDIT.
package routes

import (
	"github.com/ArcticOJ/blizzard/v0/routes/auth"
	"github.com/ArcticOJ/blizzard/v0/routes/oauth"
	"github.com/ArcticOJ/blizzard/v0/routes/posts"
	"github.com/ArcticOJ/blizzard/v0/routes/problems"
	"github.com/ArcticOJ/blizzard/v0/routes/submissions"
	"github.com/ArcticOJ/blizzard/v0/routes/user"
	"github.com/ArcticOJ/blizzard/v0/routes/users"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

var Map = map[string][]http.Route{
	"/": []http.Route{
		{
			Path:    "/status",
			Handler: GetStatus,
			Method:  http.GET,
		},
		{
			Path:    "/version",
			Handler: GetVersion,
			Method:  http.GET,
		},
	},
	"/auth": []http.Route{
		{
			Path:    "/login",
			Handler: auth.Login,
			Method:  http.POST,
		},
		{
			Path:    "/logout",
			Handler: auth.Logout,
			Method:  http.GET,
		},
		{
			Path:    "/register",
			Handler: auth.Register,
			Method:  http.POST,
		},
	},
	"/contests": []http.Route{},
	"/oauth": []http.Route{
		{
			Path:    "/",
			Handler: oauth.GetProviders,
			Method:  http.GET,
		},
		{
			Path:    "/:provider",
			Handler: oauth.CreateUrl,
			Method:  http.GET,
		},
		{
			Path:    "/:provider",
			Handler: oauth.Unlink,
			Method:  http.DELETE,
			Flags:   http.RouteAuth,
		},
		{
			Path:    "/:provider",
			Handler: oauth.Validate,
			Method:  http.POST,
		},
		{
			Path:    "/connections",
			Handler: oauth.GetConnections,
			Method:  http.GET,
			Flags:   http.RouteAuth,
		},
	},
	"/posts": []http.Route{
		{
			Path:    "/",
			Handler: posts.GetPosts,
			Method:  http.GET,
		},
	},
	"/problems": []http.Route{
		{
			Path:    "/",
			Handler: problems.GetProblems,
			Method:  http.GET,
		},
		{
			Path:    "/:id",
			Handler: problems.GetProblem,
			Method:  http.GET,
		},
		{
			Path:    "/:id/submit",
			Handler: problems.SubmitSolution,
			Method:  http.POST,
			Flags:   http.RouteAuth,
		},
	},
	"/submissions": []http.Route{
		{
			Path:    "/",
			Handler: submissions.GetSubmissions,
			Method:  http.GET,
		},
		{
			Path:    "/:id",
			Handler: submissions.GetSubmission,
			Method:  http.GET,
			Flags:   http.RouteAuth,
		},
		{
			Path:    "/:id/cancel",
			Handler: submissions.CancelSubmission,
			Method:  http.POST,
			Flags:   http.RouteAuth,
		},
		{
			Path:    "/:id/source",
			Handler: submissions.GetSourceCode,
			Method:  http.GET,
			Flags:   http.RouteAuth,
		},
	},
	"/user": []http.Route{
		{
			Path:    "/",
			Handler: user.GetUser,
			Method:  http.GET,
			Flags:   http.RouteAuth,
		},
		{
			Path:    "/:handle/hoverCard",
			Handler: user.GetHoverCard,
			Method:  http.GET,
		},
		{
			Path:    "/:handle/info",
			Handler: user.GetInfo,
			Method:  http.GET,
		},
		{
			Path:    "/:id/readme",
			Handler: user.GetReadme,
			Method:  http.GET,
		},
		{
			Path:    "/apiKey",
			Handler: user.GetApiKey,
			Method:  http.GET,
			Flags:   http.RouteAuth,
		},
		{
			Path:    "/apiKey",
			Handler: user.UpdateApiKey,
			Method:  http.PATCH,
			Flags:   http.RouteAuth,
		},
	},
	"/users": []http.Route{
		{
			Path:    "/",
			Handler: users.GetUsers,
			Method:  http.GET,
		},
	},
}
