package handler

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})

}
func (h *UserHandler) CreateNewUser(c *gin.Context) {
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()


	// Validate input
	var payload model.User


	if err := c.ShouldBindJSON(&payload); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isEmailValid(payload.Email){
		c.JSON(http.StatusInternalServerError, gin.H{"warning": "Email's format incorrect"})
		return

	}

	if len(payload.Password) > 32 || len(payload.Password) < 8{
		c.JSON(http.StatusInternalServerError, gin.H{"warning": "Your Password length must be between 8 characters and 32 characters"})
		return
	} 
	user, err := h.userService.RegisterNewUser(requestCtx, payload)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
		}
	if requestCtx.Err() != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": requestCtx.Err()})
	}
	c.JSON(http.StatusOK, gin.H{"success": user})

// Email's format
// Password Length
}

func (h *UserHandler) Login(c *gin.Context){
	var payload model.User

	if err := c.ShouldBindJSON(&payload); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Login(c, payload)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Success": user})
}

func isEmailValid(e string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(e)
}