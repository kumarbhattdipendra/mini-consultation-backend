package handlers

import (
	"fmt"
	"net/http"
	"time"

	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewAuthHandler(db *gorm.DB, secret []byte) *AuthHandler {
	return &AuthHandler{DB: db, JWTSecret: secret}
}

func (h *AuthHandler) RegisterUser(context *gin.Context) {
	var req struct{ Name, Email, Password string }
	if err := context.ShouldBindJSON(&req); err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, "Invalid JSON format.")
		return
	}

	if err := helpers.ValidateRegisterInput(req.Name, req.Email, req.Password); err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Name: req.Name, Email: req.Email, Password: string(hash)}

	fmt.Println(user)

	if err := h.DB.Create(&user).Error; err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, "Email Already exists for the user")
		return
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(h.JWTSecret)

	userData := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	response := map[string]interface{}{
		"token": signedToken,
		"user":  userData,
	}

	helpers.SendSuccessResponse(context, http.StatusCreated, response)
}

func (h *AuthHandler) LoginUser(context *gin.Context) {
	var req struct{ Email, Password string }
	if err := context.ShouldBindJSON(&req); err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, "Invalid JSON format.")
		return
	}

	if err := helpers.ValidateLoginInput(req.Email, req.Password); err != nil {
		helpers.SendErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	user, err := helpers.GetUserByEmail(h.DB, req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		helpers.SendErrorResponse(context, http.StatusUnauthorized, "invalid credentials")
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(h.JWTSecret)

	userData := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	response := map[string]interface{}{
		"token": signedToken,
		"user":  userData,
	}

	helpers.SendSuccessResponse(context, http.StatusOK, response)
}
