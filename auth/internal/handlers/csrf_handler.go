package handlers

import (
	"microservices/auth/internal/svc/csrf_svc"
	"microservices/auth/pkg/csrf_pkg"
	"time"

	"github.com/gin-gonic/gin"
)

type CSRFHandlerInterface interface {
	CsrfGet(c *gin.Context)
}

type CSRFHandlerStruct struct {
	Service csrf_svc.CsrfSvcInterface
}

func NewCSRFHandler(service csrf_svc.CsrfSvcInterface) *CSRFHandlerStruct {
	return &CSRFHandlerStruct{Service: service}
}

func (h *CSRFHandlerStruct) CsrfGet(c *gin.Context) {

	csrf := &csrf_pkg.CsrfPkgStruct{}
	token := h.Service.CreateCSRFToken(csrf, time.Now().Unix())

	c.SetCookie("csrf_token", token, 3600, "/", "", false, true)
	c.JSON(200, gin.H{
		"csrf_token": token,
	})
}
