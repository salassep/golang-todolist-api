package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/salassep/golang-todolist-api/initializers"
	"github.com/salassep/golang-todolist-api/models"
	"github.com/salassep/golang-todolist-api/utils"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Password do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Username:  payload.Username,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	userResponse := models.UserResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "username = ?", strings.ToLower(payload.Username))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid username or password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid username or password"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	acessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", acessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": acessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Could not refresh access token"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "The user belonging to this token no longer exists"})
		return
	}

	acessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", acessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": acessToken})
}

func (as *AuthController) LogOutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
