package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/insanjati/fitbyte/internal/service"
)

type AuthMiddleware interface {
	CheckToken() gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

func (a *authMiddleware) CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing request token",
			})
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid request token format",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing request token",
			})
			return
		}

		claims, err := a.jwtService.VerifyToken(token)
		if err != nil {

			var message string
			switch {
			case strings.Contains(err.Error(), "expired"):
				message = "Expired request token"
			case strings.Contains(err.Error(), "invalid"):
				message = "Invalid request token"
			default:
				message = "Invalid request token"
			}

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": message,
			})
			return
		}

		uidStr, ok := claims["user_id"].(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid user ID type",
			})
			return
		}

		uid, err := uuid.Parse(uidStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid UUID format",
			})
			return
		}

		ctx.Set("user_id", uid)
		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
