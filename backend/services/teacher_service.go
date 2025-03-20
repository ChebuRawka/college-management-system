package services

import (
	"backend/models"
	"backend/repository"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type TeacherService struct {
    Repo *repositories.TeacherRepository
}

func NewTeacherService(repo *repositories.TeacherRepository) *TeacherService {
    return &TeacherService{Repo: repo}
}

// Создание преподавателя
func (s *TeacherService) CreateTeacher(teacher *models.Teacher) error {
    // Проверяем валидацию модели
    if err := models.Validate.Struct(teacher); err != nil {
        return err
    }

    // Проверяем, существует ли уже преподаватель с таким именем и предметом
    exists, err := s.Repo.CheckTeacherExists(teacher.Name, teacher.Subject)
    if err != nil {
        return err
    }
    if exists {
        return errors.New("teacher with this name and subject already exists")
    }

    // Проверяем, что все указанные курсы существуют
    if len(teacher.Courses) > 0 {
        validCourses, err := s.Repo.CheckCoursesExist(teacher.Courses)
        if err != nil {
            return err
        }
        if !validCourses {
            return errors.New("some courses do not exist")
        }
    }

    // Создаем нового преподавателя
    return s.Repo.CreateTeacher(teacher)
}

// Получение всех преподавателей
func (s *TeacherService) GetAllTeachers() ([]models.Teacher, error) {
    return s.Repo.GetAllTeachers()
}

// Получение преподавателя по ID
func (s *TeacherService) GetTeacherByID(id int) (*models.Teacher, error) {
    return s.Repo.GetTeacherByID(id)
}

func (s *TeacherService) UpdateTeacherPartial(teacherID int, updates map[string]interface{}) error {
    err := s.Repo.UpdateTeacherPartial(teacherID, updates)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            return fmt.Errorf("teacher with idiot %d not found", teacherID)
        }
        return fmt.Errorf("failed to update teacher: %v", err)
    }
    return nil
}

// Удаление преподавателя
func (s *TeacherService) DeleteTeacher(id int) error {
    return s.Repo.DeleteTeacher(id)
}

func (s *TeacherService) GetAllTeachersWithCourses() ([]models.Teacher, error) {
    return s.Repo.GetAllTeachersWithCourses()
}

func (s *TeacherService) GetTeacherSchedule(teacherName string) ([]models.ScheduleResponse, error) {
    return s.Repo.GetTeacherSchedule(teacherName)
}

func (s *TeacherService) UpdateTeacherProfile(teacherID int, updates map[string]interface{}) error {
    // Если передан новый пароль, хэшируем его
    if newPassword, ok := updates["password"].(string); ok {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        updates["password"] = string(hashedPassword)
    }

    // Обновляем данные в базе данных
    if err := s.Repo.UpdateTeacherProfile(teacherID, updates); err != nil {
        return err
    }

    return nil
}