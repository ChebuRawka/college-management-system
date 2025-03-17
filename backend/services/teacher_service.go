package services

import (
    "backend/models"
    "backend/repository"
)

type TeacherService struct {
    Repo *repositories.TeacherRepository
}

func NewTeacherService(repo *repositories.TeacherRepository) *TeacherService {
    return &TeacherService{Repo: repo}
}

// Создание преподавателя
func (s *TeacherService) CreateTeacher(teacher models.Teacher) (models.Teacher, error) {
    id, err := s.Repo.CreateTeacher(teacher)
    if err != nil {
        return models.Teacher{}, err
    }
    teacher.ID = id
    return teacher, nil
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