package handlers

import (
	"microservices/auth/internal/models"
	"microservices/auth/internal/svc/clock_svc"
	"microservices/auth/internal/svc/jwt_svc"
	"microservices/auth/pkg/encrypt_pkg"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterHandlerInterface interface {
	HandleRegister(c *gin.Context)
}

type RegisterHandlerStruct struct {
	Db          *gorm.DB
	jwt_svc     jwt_svc.JwtServiceInterface
	encrypt_pkg encrypt_pkg.EncryptPkgInterface
	Clock       clock_svc.ClockInterface
}

func NewRegisterHandler(
	db *gorm.DB,
	jwt_svc jwt_svc.JwtServiceInterface,
	encrypt_pkg encrypt_pkg.EncryptPkgInterface,
	clock clock_svc.ClockInterface,
) *RegisterHandlerStruct {
	return &RegisterHandlerStruct{
		Db:          db,
		jwt_svc:     jwt_svc,
		encrypt_pkg: encrypt_pkg,
		Clock:       clock,
	}
}

type registerRequest struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}

func (h *RegisterHandlerStruct) HandleRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	hashedPassword, err := h.encrypt_pkg.CreatePasswordHash(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     string(hashedPassword),
		RefreshToken: h.jwt_svc.CreateRefreshToken(h.Clock),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := h.Db.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully", "user": user})
}
