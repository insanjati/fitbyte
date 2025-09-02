package handler

import (
	"net/http"

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
// Validate input
var payload model.User
if err := c.ShouldBindJSON(&payload); err != nil{
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + "Test2"})
	return
}



data, err := h.userService.RegisterNewUser(c, payload)
	if err != nil{
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + "Test1"})
	return
}
		c.JSON(http.StatusOK, gin.H{"success": data})

// Email's format
// Password Length
}
