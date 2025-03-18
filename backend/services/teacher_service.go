package services

import (
    "backend/models"
    "backend/repository"
    "errors"
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

// Частичное обновление преподавателя
func (s *TeacherService) UpdateTeacherPartial(id int, updates map[string]interface{}) (*models.Teacher, error) {
    updatedTeacher, err := s.Repo.UpdateTeacherPartial(id, updates)
    if err != nil {
        return nil, err
    }
    return updatedTeacher, nil
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