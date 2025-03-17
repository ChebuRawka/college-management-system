package services

import (
    "backend/models"
    "backend/repository"
)

type StudentService struct {
    Repo *repositories.StudentRepository
}

func NewStudentService(repo *repositories.StudentRepository) *StudentService {
    return &StudentService{Repo: repo}
}

func (s *StudentService) CreateStudent(student *models.Student) error {
    return s.Repo.CreateStudent(student)
}

func (s *StudentService) GetStudents() ([]models.Student, error) {
    return s.Repo.GetStudents()
}

func (s *StudentService) GetStudentByID(id int) (*models.Student, error) {
    return s.Repo.GetStudentByID(id)
}

func (s *StudentService) UpdateStudent(id int, updates map[string]interface{}) (*models.Student, error) {
    return s.Repo.UpdateStudent(id, updates)
}

func (s *StudentService) DeleteStudent(id int) error {
    return s.Repo.DeleteStudent(id)
}

func (s *StudentService) CourseExists(courseName string) (bool, error) {
    return s.Repo.CourseExists(courseName)
}