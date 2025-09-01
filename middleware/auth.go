package middleware

import (
	"backend/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.SendErrorResponse(context, http.StatusUnauthorized, "Authentication Required.")
			context.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			helpers.SendErrorResponse(context, http.StatusUnauthorized, "Invalid token")
			context.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if claims != nil {
			context.Set("userId", uint(claims["user_id"].(float64)))
		}

		context.Next()
	}
}
