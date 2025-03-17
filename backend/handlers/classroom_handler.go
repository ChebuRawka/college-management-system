package handlers

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassroomHandler struct {
    Service *services.ClassroomService
}

func NewClassroomHandler(service *services.ClassroomService) *ClassroomHandler {
    return &ClassroomHandler{Service: service}
}

func (h *ClassroomHandler) CreateClassroom(c *gin.Context) {
    var classroom models.Classroom
    if err := c.ShouldBindJSON(&classroom); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if err := h.Service.CreateClassroom(&classroom); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, classroom)
}

func (h *ClassroomHandler) GetClassrooms(c *gin.Context) {
    classrooms, err := h.Service.GetClassrooms()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, classrooms)
}

func (h *ClassroomHandler) GetClassroomByID(c *gin.Context) {
    id := c.Param("id")
    classroomID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    classroom, err := h.Service.GetClassroomByID(classroomID)
    if err != nil {
        if errors.Is(err, errors.New("classroom not found")) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Classroom not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, classroom)
}

func (h *ClassroomHandler) UpdateClassroom(c *gin.Context) {
    id := c.Param("id")
    classroomID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var updates map[string]interface{}
    if err := json.NewDecoder(c.Request.Body).Decode(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    classroom, err := h.Service.UpdateClassroom(classroomID, updates)
    if err != nil {
        if errors.Is(err, errors.New("classroom not found")) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Classroom not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, classroom)
}

func (h *ClassroomHandler) DeleteClassroom(c *gin.Context) {
    id := c.Param("id")
    classroomID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := h.Service.DeleteClassroom(classroomID); err != nil {
        if errors.Is(err, errors.New("classroom not found")) {
            c.JSON(http.StatusNotFound, gin.H{"error": "Classroom not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Classroom deleted successfully"})
}