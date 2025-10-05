package app

import (
	"microservices/chat/internal/handlers"
	"microservices/chat/internal/middlewares"
	"microservices/chat/internal/routings"
	"microservices/chat/internal/svc/chat_svc"
	"microservices/chat/internal/svc/clock_svc"
	"microservices/chat/internal/svc/csrf_svc"
	"microservices/chat/internal/svc/mongo_svc"
	"microservices/chat/pkg/csrf_pkg"
	"microservices/chat/pkg/mongo_pkg"

	"github.com/gin-gonic/gin"
)

type App struct {
	CsrfMW   gin.HandlerFunc
	AuthMW   gin.HandlerFunc
	Handlers *handlers.HandlerStruct
}

func NewApp() *App {
	csrfPkg := &csrf_pkg.CsrfPkgStruct{}
	mongoPkg := mongo_pkg.NewMongoPkg()

	verifier := csrf_svc.NewVerifier(csrfPkg, "secrets", clock_svc.RealClockStruct{})

	csrfMW := middlewares.NewCSRFMiddleware(verifier)
	authMW := middlewares.NewAuthMiddleware()

	mongoSvc := mongo_svc.NewMongoSvc(&mongo_pkg.RealMongoDatabase{})

	chatSvc := chat_svc.NewChatSvc()

	app := &App{
		CsrfMW:   csrfMW.Handler(),
		AuthMW:   authMW.Handler(),
		Handlers: handlers.NewHandlers(mongoSvc, mongoPkg, chatSvc),
	}
	return app
}

func (a *App) InitRoutes(r *gin.Engine) {
	routings.Routing(r, a.CsrfMW, a.AuthMW, a.Handlers)
}
