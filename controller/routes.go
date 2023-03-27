package controller

import (
	"net/http"

	"github.com/gkontos/goapi/logger"
	"github.com/gkontos/goapi/security"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/net/context"
)

func (api *apiRouter) SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&StructuredLogger{&logger.Logger}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// add Cors Headers
	r.Use(api.rs.CorsHeaders)

	r.Route("/v1", func(r chi.Router) {
		r.Use(apiVersionCtx("v1"))
		r.Mount("/users", userRouter(api.rs, api.ctrl))
		r.Mount("/login", tokenRouter(api.ctrl))
	})

	return r
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			next.ServeHTTP(w, r)
		})
	}
}

// use authenticate and authorize middleware
func userRouter(rs security.RouterSecurity, ctrl *apiController) chi.Router {
	r := chi.NewRouter()

	r.Use(rs.AuthenticateAuthHeader)
	r.Get("/",
		AddMiddleware(
			http.HandlerFunc(ctrl.GetUsers),
			rs.Authorize(security.Permission("admin"))))
	return r
}

// use no auth middleware
func tokenRouter(ctrl *apiController) chi.Router {
	r := chi.NewRouter()
	r.Post("/refresh", ctrl.TokenRefresh)
	r.Post("/", ctrl.TokenCreate)
	return r
}

// AddMiddleware will add functions before processing the main request
// NOTE : The middleware functions run in reverse order ... at least functionally.  eg If you want to authenticate a
// request and then authorize the principle, load the middleware functions in the order authorize, authenticate
func AddMiddleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}
