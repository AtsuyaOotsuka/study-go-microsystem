package handlers

import (
	"microservices/auth/internal/jwt_svc"
	"microservices/auth/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandlerInterface interface {
	HandleLogin(c *gin.Context)
	HandleRefresh(c *gin.Context)
}

type AuthHandlerStruct struct {
	Db      *gorm.DB
	jwt_svc jwt_svc.JwtServiceInterface
}

func NewAuthHandler(db *gorm.DB, jwtSvc jwt_svc.JwtServiceInterface) *AuthHandlerStruct {
	return &AuthHandlerStruct{
		Db:      db,
		jwt_svc: jwtSvc,
	}
}

type loginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}

func (h *AuthHandlerStruct) HandleLogin(c *gin.Context) {
	// 1) 入力バインド（JSONでもx-www-form-urlencodedでもOK）
	var req loginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// 2) 必要カラムだけ取得（ID & Password）
	userRepository := repositories.UserRepositoryStruct{Db: h.Db}
	user, err := userRepository.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"}) // 本来は曖昧にするが、学習目的のため、分ける
		return
	}

	// 3) パスワード検証（models.User#VerifyPassword が bcrypt.CompareHashAndPassword を呼ぶ想定）
	if err := user.VerifyPassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"}) // 本来は曖昧にするが、学習目的のため、分ける
		return
	}

	// JWTトークンを作成
	tokenString, err := h.jwt_svc.CreateJwt(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT"})
		return
	}

	resp := map[string]interface{}{
		"access_token":  tokenString,
		"refresh_token": user.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}

	c.JSON(http.StatusOK, resp)
}

type refreshRequest struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

func (h *AuthHandlerStruct) HandleRefresh(c *gin.Context) {

	var req refreshRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userRepository := repositories.UserRepositoryStruct{Db: h.Db}

	user, err := userRepository.GetByRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// JWTトークンを作成
	tokenString, err := h.jwt_svc.CreateJwt(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT"})
		return
	}

	resp := map[string]interface{}{
		"access_token":  tokenString,
		"refresh_token": user.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}

	c.JSON(http.StatusOK, resp)
}
