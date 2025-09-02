package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		// Check if Authorization header exists
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

		// Set user context
		ctx.Set("user_id", claims["user_id"])
		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
