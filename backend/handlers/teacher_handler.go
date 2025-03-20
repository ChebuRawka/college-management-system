package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"backend/models"
	"backend/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
    Service *services.TeacherService
    EmailService *services.EmailService
}

func NewTeacherHandler(service *services.TeacherService, emailService *services.EmailService) *TeacherHandler {
    return &TeacherHandler{
        Service:     service,
        EmailService: emailService,
    }
}

// Отправка уведомления преподавателю
func (h *TeacherHandler) NotifyTeacher(c *gin.Context) {
    var input struct {
        Email   string `json:"email"`
        Subject string `json:"subject"`
        Body    string `json:"body"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if err := h.EmailService.SendEmail(input.Email, input.Subject, input.Body); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
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
    // Получаем ID преподавателя из параметров URL
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    // Привязываем JSON-данные к словарю updates
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if len(updates) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
        return
    }

    // Вызываем метод сервиса для обновления данных
    err = h.Service.UpdateTeacherPartial(id, updates)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Teacher updated successfully"})
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

func (h *TeacherHandler) UpdateTeacherProfile(c *gin.Context) {
    teacherID, exists := c.Get("user_id") // Получаем ID пользователя из контекста (JWT)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        fmt.Println("Error binding JSON:", err) // Логируем ошибку привязки JSON
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    fmt.Println("Received updates:", updates) // Логируем входные данные

    if len(updates) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
        return
    }

    if err := h.Service.UpdateTeacherProfile(teacherID.(int), updates); err != nil {
        fmt.Println("Service error:", err) // Логируем ошибку сервиса
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}