package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "backend/models"
    "backend/services"
    "backend/utils"
	"fmt"
	"time"
)

type StudentHandler struct {
    Service *services.StudentService
}

func NewStudentHandler(service *services.StudentService) *StudentHandler {
    return &StudentHandler{Service: service}
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
    var student models.Student

    // Получаем данные из запроса
    if err := c.ShouldBindJSON(&student); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Проверяем существование курса с указанным group_name
    exists, err := h.Service.CourseExists(student.GroupName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("course with name '%s' does not exist", student.GroupName)})
        return
    }

    // Проверяем и парсим дату рождения
    dateOfBirth, err := time.Parse("2006-01-02", student.DateOfBirth)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_of_birth format. Use YYYY-MM-DD"})
        return
    }

    // Вычисляем возраст
    student.Age = utils.CalculateAge(dateOfBirth)

    // Передаём дату рождения в сервис
    student.DateOfBirth = dateOfBirth.Format("2006-01-02") // Сохраняем в формате YYYY-MM-DD

    // Создаём студента
    if err := h.Service.CreateStudent(&student); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, student)
}
func (h *StudentHandler) GetStudents(c *gin.Context) {
    students, err := h.Service.GetStudents()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    student, err := h.Service.GetStudentByID(id)
    if err != nil {
        if err.Error() == fmt.Sprintf("student with id %d not found", id) {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
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

    // Если есть обновление group_name, проверяем его существование
    if groupName, ok := updates["group_name"].(string); ok {
        exists, err := h.Service.CourseExists(groupName)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if !exists {
            c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("course with name '%s' does not exist", groupName)})
            return
        }
    }

    // Если есть обновление даты рождения
    if dateOfBirth, ok := updates["date_of_birth"].(string); ok {
        parsedDate, err := time.Parse("2006-01-02", dateOfBirth)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_of_birth format. Use YYYY-MM-DD"})
            return
        }
        updates["date_of_birth"] = parsedDate.Format("2006-01-02") // Сохраняем в формате YYYY-MM-DD
    }

    updatedStudent, err := h.Service.UpdateStudent(id, updates)
    if err != nil {
        if err.Error() == fmt.Sprintf("student with id %d not found", id) {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        if err.Error() == "no fields to update" {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Вычисляем возраст
    dateOfBirth, _ := time.Parse("2006-01-02", updatedStudent.DateOfBirth)
    updatedStudent.Age = utils.CalculateAge(dateOfBirth)

    c.JSON(http.StatusOK, updatedStudent)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := h.Service.DeleteStudent(id); err != nil {
        if err.Error() == fmt.Sprintf("student with id %d not found", id) {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}