package services

import (
    "backend/models"
    "backend/repository"
)

type ClassroomService struct {
    Repo *repositories.ClassroomRepository
}

func NewClassroomService(repo *repositories.ClassroomRepository) *ClassroomService {
    return &ClassroomService{Repo: repo}
}

func (s *ClassroomService) CreateClassroom(classroom *models.Classroom) error {
    return s.Repo.CreateClassroom(classroom)
}

func (s *ClassroomService) GetClassrooms() ([]models.Classroom, error) {
    return s.Repo.GetClassrooms()
}

func (s *ClassroomService) GetClassroomByID(id int) (*models.Classroom, error) {
    return s.Repo.GetClassroomByID(id)
}

func (s *ClassroomService) UpdateClassroom(id int, updates map[string]interface{}) (*models.Classroom, error) {
    return s.Repo.UpdateClassroom(id, updates)
}

func (s *ClassroomService) DeleteClassroom(id int) error {
    return s.Repo.DeleteClassroom(id)
}