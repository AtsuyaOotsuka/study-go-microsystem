package app

import (
	"database/sql"
	"microservices/auth/internal/handlers"
	"microservices/auth/internal/middlewares"
	"microservices/auth/internal/routings"
	"microservices/auth/internal/svc/clock_svc"
	"microservices/auth/internal/svc/csrf_svc"
	"microservices/auth/internal/svc/jwt_svc"
	"microservices/auth/pkg/csrf_pkg"
	"microservices/auth/pkg/encrypt_pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	CSRFHandler        *handlers.CSRFHandlerStruct
	AuthHandler        *handlers.AuthHandlerStruct
	HealthCheckHandler *handlers.HealthCheckHandlerStruct
	RegisterHandler    *handlers.RegisterHandlerStruct

	CsrfMW gin.HandlerFunc
}

func NewApp(db *gorm.DB, sqlDB *sql.DB) (*App, func(), error) {

	csrf_pkg := &csrf_pkg.CsrfPkgStruct{}
	encrypt_pkg := &encrypt_pkg.EncryptPkgStruct{}

	verifier := csrf_svc.NewVerifier(csrf_pkg, "secrets", clock_svc.RealClockStruct{})
	csrfMW := middlewares.NewCSRFMiddleware(verifier)

	jwtSvc := jwt_svc.NewJwtService()

	app := &App{
		CSRFHandler:        handlers.NewCSRFHandler(&csrf_svc.CsrfSvcStruct{}),
		AuthHandler:        handlers.NewAuthHandler(db, jwtSvc),
		HealthCheckHandler: handlers.NewHealthCheckHandler(),
		RegisterHandler:    handlers.NewRegisterHandler(db, jwtSvc, encrypt_pkg, clock_svc.RealClockStruct{}),
		CsrfMW:             csrfMW.Handler(),
	}

	cleanup := func() { sqlDB.Close() }
	return app, cleanup, nil
}

func (a *App) InitRoutes(r *gin.Engine) {
	routings.CsrfRouting(r, a.CSRFHandler)
	routings.HealthCheckRouting(r, a.HealthCheckHandler)
	routings.AuthRouting(r, a.AuthHandler, a.CsrfMW)
	routings.RegisterRouting(r, a.RegisterHandler, a.CsrfMW)
}
