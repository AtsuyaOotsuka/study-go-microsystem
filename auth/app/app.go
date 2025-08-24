package app

import (
	"database/sql"
	"microservices/auth/handlers"
	"microservices/auth/internal/clock_svc"
	"microservices/auth/internal/csrf_svc"
	"microservices/auth/internal/jwt_svc"
	"microservices/auth/middlewares"
	"microservices/auth/pkg/csrf_pkg"
	"microservices/auth/routings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	CSRFHandler        *handlers.CSRFHandlerStruct
	AuthHandler        *handlers.AuthHandlerStruct
	HealthCheckHandler *handlers.HealthCheckHandlerStruct

	CsrfMW gin.HandlerFunc
}

func NewApp(db *gorm.DB, sqlDB *sql.DB) (*App, func(), error) {

	csrf_pkg := &csrf_pkg.CsrfPkgStruct{}
	verifier := csrf_svc.NewVerifier(csrf_pkg, "secrets", clock_svc.RealClockStruct{})
	csrfMW := middlewares.NewCSRFMiddleware(verifier)

	jwtSvc := jwt_svc.NewJwtService()

	app := &App{
		CSRFHandler:        handlers.NewCSRFHandler(&csrf_svc.CsrfSvcStruct{}),
		AuthHandler:        handlers.NewAuthHandler(db, jwtSvc),
		HealthCheckHandler: handlers.NewHealthCheckHandler(),
		CsrfMW:             csrfMW.Handler(),
	}

	cleanup := func() { sqlDB.Close() }
	return app, cleanup, nil
}

func (a *App) InitRoutes(r *gin.Engine) {
	routings.CsrfRouting(r, a.CSRFHandler)
	healthCheckHandler := handlers.NewHealthCheckHandler()
	routings.HealthCheckRouting(r, healthCheckHandler)
	routings.AuthRouting(r, a.AuthHandler, a.CsrfMW)
}
