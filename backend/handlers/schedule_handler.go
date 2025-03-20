package handlers

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
    Service *services.ScheduleService
}

func NewScheduleHandler(service *services.ScheduleService) *ScheduleHandler {
    return &ScheduleHandler{Service: service}
}

// CreateSchedule создает новую запись в расписании
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
    type RequestBody struct {
        TeacherID   int       `json:"teacher_id"`
        ClassroomID int       `json:"classroom_id"`
        GroupName   string    `json:"group_name"`
        StartTime   time.Time `json:"start_time"`
        EndTime     time.Time `json:"end_time"`
        DayOfWeek   string    `json:"day_of_week"`
    }

    var req RequestBody
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    schedule := &models.Schedule{
        GroupName: req.GroupName,
        StartTime: req.StartTime,
        EndTime:   req.EndTime,
        DayOfWeek: req.DayOfWeek,
    }

    if err := h.Service.CreateSchedule(req.TeacherID, req.ClassroomID, schedule); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
    schedules, err := h.Service.GetSchedules()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, schedules)
}

// GetScheduleByID возвращает запись расписания по ID
func (h *ScheduleHandler) GetScheduleByID(c *gin.Context) {
    id := c.Param("id")
    scheduleID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    schedule, err := h.Service.GetScheduleByID(scheduleID)
    if err != nil {
        if errors.Is(err, errors.New("schedule not found")) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, schedule)
}

// UpdateSchedule обновляет запись расписания
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
    id := c.Param("id")
    scheduleID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var updates map[string]interface{}
    if err := json.NewDecoder(c.Request.Body).Decode(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    schedule, err := h.Service.UpdateSchedule(scheduleID, updates)
    if err != nil {
        if errors.Is(err, errors.New("schedule not found")) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, schedule)
}

// DeleteSchedule удаляет запись расписания по ID
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
    id := c.Param("id")
    scheduleID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := h.Service.DeleteSchedule(scheduleID); err != nil {
        if errors.Is(err, errors.New("schedule not found")) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}

func (h *ScheduleHandler) GetSchedulesByDay(c *gin.Context) {
    dayOfWeek := c.Param("day")
    if dayOfWeek == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "day_of_week is required"})
        return
    }

    schedules, err := h.Service.GetSchedulesByDay(dayOfWeek)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, schedules)
}

// GetSchedulesByGroup возвращает расписание для конкретной группы
func (h *ScheduleHandler) GetSchedulesByGroup(c *gin.Context) {
    groupName := c.Param("group_name")
    
    // Логируем полученное значение group_name
    fmt.Printf("Received group_name: %q\n", groupName)

    if groupName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Group name is required"})
        return
    }

    schedules, err := h.Service.GetSchedulesByGroup(groupName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, schedules)
}


