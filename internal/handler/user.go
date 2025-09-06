package handler

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c *gin.Context) {

	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := uid.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	users, err := h.userService.FindUserById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)

}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			// Custom error messages
			errorMessages := make([]string, 0, len(errs))
			for _, e := range errs {
				switch e.Field() {
				case "Preference":
					if e.Tag() == "oneof" {
						errorMessages = append(errorMessages, "preference must be CARDIO or WEIGHT")
					} else {
						errorMessages = append(errorMessages, "preference is required")
					}
				case "WeightUnit":
					errorMessages = append(errorMessages, "weightUnit must be KG or LBS")
				case "HeightUnit":
					errorMessages = append(errorMessages, "heightUnit must be CM or INCH")
				case "Weight":
					errorMessages = append(errorMessages, "weight must be between 10 and 1000")
				case "Height":
					errorMessages = append(errorMessages, "height must be between 3 and 250")
				default:
					errorMessages = append(errorMessages, fmt.Sprintf("%s is invalid", e.Field()))
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
	}

	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := uid.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}
	updatedUser, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) CreateNewUser(c *gin.Context) {
	requestCtx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Validate input
	var payload model.User

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isEmailValid(payload.Email) {
		c.JSON(http.StatusInternalServerError, gin.H{"warning": "Email's format incorrect"})
		return

	}

	if len(payload.Password) > 32 || len(payload.Password) < 8 {
		c.JSON(http.StatusInternalServerError, gin.H{"warning": "Your Password length must be between 8 characters and 32 characters"})
		return
	}
	user, err := h.userService.RegisterNewUser(requestCtx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if requestCtx.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": requestCtx.Err()})
	}
	c.JSON(http.StatusOK, gin.H{"success": user})

	// Email's format
	// Password Length
}

func (h *UserHandler) Login(c *gin.Context) {
	var payload model.User

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Login(c, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Success": user})
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
