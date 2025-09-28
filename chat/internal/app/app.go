package app

import (
	"microservices/chat/internal/handlers"
	"microservices/chat/internal/middlewares"
	"microservices/chat/internal/routings"
	"microservices/chat/internal/svc/clock_svc"
	"microservices/chat/internal/svc/csrf_svc"
	"microservices/chat/internal/svc/mongo_svc"
	"microservices/chat/pkg/csrf_pkg"

	"github.com/gin-gonic/gin"
)

type App struct {
	CsrfMW   gin.HandlerFunc
	AuthMW   gin.HandlerFunc
	Handlers *handlers.HandlerStruct
}

func NewApp() *App {
	csrf_pkg := &csrf_pkg.CsrfPkgStruct{}

	verifier := csrf_svc.NewVerifier(csrf_pkg, "secrets", clock_svc.RealClockStruct{})

	csrfMW := middlewares.NewCSRFMiddleware(verifier)
	authMW := middlewares.NewAuthMiddleware()

	mongoSvc := mongo_svc.NewMongoSvc()

	app := &App{
		CsrfMW:   csrfMW.Handler(),
		AuthMW:   authMW.Handler(),
		Handlers: handlers.NewHandlers(mongoSvc),
	}
	return app
}

func (a *App) InitRoutes(r *gin.Engine) {
	routings.Routing(r, a.CsrfMW, a.AuthMW, a.Handlers)
}
