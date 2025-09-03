package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/insanjati/fitbyte/internal/middleware"
	"github.com/insanjati/fitbyte/internal/model"
	"github.com/insanjati/fitbyte/internal/service"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	activityService *service.ActivityService
}

func NewActivityHandler(activityService *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

func getUserID(c *gin.Context) int {
	val, _ := c.Get(middleware.ContextUserIDKey)
	if id, ok := val.(int); ok {
		return id
	}
	return 0
}

// POST /v1/activity
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var req model.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error"})
		return
	}

	userID := getUserID(c)

	activity, err := h.activityService.CreateActivity(userID, req)
	fmt.Print(err)

	if err != nil {
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
	userID := getUserID(c)

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
