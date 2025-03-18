package handlers

import (
    "net/http"
    "strconv"

    "backend/models"
    "backend/services"
    "github.com/gin-gonic/gin"
    "fmt"
    "strings"
)

type TeacherHandler struct {
    Service *services.TeacherService
}

func NewTeacherHandler(service *services.TeacherService) *TeacherHandler {
    return &TeacherHandler{Service: service}
}

// Создание преподавателя
func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
    var teacher models.Teacher
    if err := c.ShouldBindJSON(&teacher); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    fmt.Printf("Received teacher data: %+v\n", teacher)

    // Пытаемся создать преподавателя
    if err := h.Service.CreateTeacher(&teacher); err != nil {
        if err.Error() == "teacher with this name and subject already exists" {
            c.JSON(http.StatusConflict, gin.H{"error": "Teacher with this name and subject already exists"})
            return
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, teacher)
}

// Получение всех преподавателей
func (h *TeacherHandler) GetAllTeachers(c *gin.Context) {
    teachers, err := h.Service.GetAllTeachers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    fmt.Printf("Returning teachers: %+v\n", teachers)

    c.JSON(http.StatusOK, teachers)
}



// Частичное обновление преподавателя
func (h *TeacherHandler) UpdateTeacherPartial(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    updatedTeacher, err := h.Service.UpdateTeacherPartial(id, updates)
    if err != nil {
        if err.Error() == fmt.Sprintf("teacher with id %d not found", id) {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        if err.Error() == "no fields to update" {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if strings.Contains(err.Error(), "invalid field") {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedTeacher)
}

// Удаление преподавателя
func (h *TeacherHandler) DeleteTeacher(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    err = h.Service.DeleteTeacher(id)
    if err != nil {
        if err.Error() == fmt.Sprintf("teacher with id %d not found", id) {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Teacher deleted successfully"})
}

func (h *TeacherHandler) GetTeacherSchedule(c *gin.Context) {
    teacherName := c.Param("teacher_name")
    fmt.Println("Fetching schedule for teacher:", teacherName)

    if teacherName == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "teacher_name is required"})
        return
    }

    schedules, err := h.Service.GetTeacherSchedule(teacherName)
    if err != nil {
        fmt.Println("Error fetching teacher schedule:", err)
        if err.Error() == "teacher not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    fmt.Println("Fetched schedule:", schedules)
    c.JSON(http.StatusOK, schedules)
}