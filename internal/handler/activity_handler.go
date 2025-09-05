package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/insanjati/fitbyte/internal/errors"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivityHandler struct {
	activityService *service.ActivityService
}

func NewActivityHandler(activityService *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

func getUserID(c *gin.Context) (uuid.UUID, error) {
	rawUserID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.ErrUnauthorized
	}
	userID, ok := rawUserID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("Invalid user ID type")
	}
	return userID, nil
}

// POST /v1/activity
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req model.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	userID, err := getUserID(c)
	if err == errors.ErrUnauthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	activity, err := h.activityService.CreateActivity(userID, req)

	if err != nil {
		fmt.Printf("CreateActivity error: %v\n", err)
		if err == errors.ErrUnauthorized {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrUnauthorized.Error()})
			return
		}
		if err.Error() == "invalid doneAt" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doneAt"})
			return
		}
		if err.Error() == "invalid activityType" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activityType"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"activityId":        activity.ID,
		"activityType":      activity.ActivityType,
		"doneAt":            activity.DoneAt.Format(time.RFC3339),
		"durationInMinutes": activity.DurationInMinutes,
		"caloriesBurned":    activity.CaloriesBurned,
		"createdAt":         activity.CreatedAt.Format(time.RFC3339),
		"updatedAt":         activity.UpdatedAt.Format(time.RFC3339),
	})
}

func (h *ActivityHandler) GetUserActivities(c *gin.Context) {
	userID, err := getUserID(c)
	if err == errors.ErrUnauthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	var filter model.ActivityFilter

	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filter.Limit = &n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			filter.Offset = &n
		}
	}
	if v := c.Query("activityType"); v != "" {
		at := model.ActivityType(v)
		if _, ok := model.ActivityTypeCalories[at]; ok {
			filter.ActivityType = &at
		}
	}
	if v := c.Query("doneAtFrom"); v != "" {
		if _, err := time.Parse(time.RFC3339, v); err == nil {
			filter.DoneAtFrom = &v
		}
	}
	if v := c.Query("doneAtTo"); v != "" {
		if _, err := time.Parse(time.RFC3339, v); err == nil {
			filter.DoneAtTo = &v
		}
	}
	if v := c.Query("caloriesBurnedMin"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.CaloriesBurnedMin = &n
		}
	}
	if v := c.Query("caloriesBurnedMax"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.CaloriesBurnedMax = &n
		}
	}

	activities, err := h.activityService.GetUserActivities(userID, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	resp := make([]gin.H, 0, len(activities))
	for _, a := range activities {
		resp = append(resp, gin.H{
			"activityId":        a.ID,
			"activityType":      a.ActivityType,
			"doneAt":            a.DoneAt.Format(time.RFC3339),
			"durationInMinutes": a.DurationInMinutes,
			"caloriesBurned":    a.CaloriesBurned,
			"createdAt":         a.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, resp)
}

// PATCH /v1/activity
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	var req model.UpdateActivityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	userID, err := getUserID(c)
	if err == errors.ErrUnauthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	// Validate activityID
	activityID, err := uuid.Parse(c.Param("activityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "invalid activityId",
			"description": "activityId must be a valid UUID",
		})
		return
	}

	activity, err := h.activityService.UpdateActivity(userID, activityID, req)
	fmt.Print(err)

	if err != nil {
		if err == errors.ErrUnauthorized {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrUnauthorized.Error()})
			return
		}
		if err.Error() == "invalid doneAt" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doneAt"})
			return
		}
		if err.Error() == "invalid activityType" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activityType"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"activityId":        activity.ID,
		"activityType":      activity.ActivityType,
		"doneAt":            activity.DoneAt.Format(time.RFC3339),
		"durationInMinutes": activity.DurationInMinutes,
		"caloriesBurned":    activity.CaloriesBurned,
		"createdAt":         activity.CreatedAt.Format(time.RFC3339),
		"updatedAt":         activity.UpdatedAt.Format(time.RFC3339),
	})
}

// DELETE /v1/activity/:activityId
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	// Get activity ID from URL parameter
	activityIDStr := c.Param("activityId")
	activityID, err := uuid.Parse(activityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity ID format"})
		return
	}

	// Get user ID from JWT context
	userID, err := getUserID(c)
	if err == errors.ErrUnauthorized {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	// Delete the activity
	err = h.activityService.DeleteActivity(activityID, userID)
	if err != nil {
		if err.Error() == "activity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
