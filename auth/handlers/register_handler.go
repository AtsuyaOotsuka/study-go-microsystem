package handlers

import (
	"microservices/auth/internal/clock_svc"
	"microservices/auth/internal/jwt_svc"
	"microservices/auth/models"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterHandlerInterface interface {
	Register(c *gin.Context)
}

type RegisterHandlerStruct struct {
	Db      *gorm.DB
	jwt_svc jwt_svc.JwtServiceInterface
	Clock   clock_svc.ClockInterface
}

func NewRegisterHandler(db *gorm.DB, jwt_svc jwt_svc.JwtServiceInterface, clock clock_svc.ClockInterface) *RegisterHandlerStruct {
	return &RegisterHandlerStruct{
		Db:      db,
		jwt_svc: jwt_svc,
		Clock:   clock,
	}
}

type registerRequest struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}

func (h *RegisterHandlerStruct) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     string(hashedPassword),
		RefreshToken: h.jwt_svc.CreateRefreshToken(h.Clock),
	}

	result := h.Db.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully", "user": user})
}
