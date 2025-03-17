package services

import (
    "backend/models"
    "backend/repository"
)

type CourseService struct {
    Repo *repositories.CourseRepository
}

func NewCourseService(repo *repositories.CourseRepository) *CourseService {
    return &CourseService{Repo: repo}
}

func (s *CourseService) CreateCourse(course *models.Course) error {
    return s.Repo.CreateCourse(course)
}

func (s *CourseService) GetCourses() ([]models.Course, error) {
    return s.Repo.GetCourses()
}

func (s *CourseService) GetCourseByID(id int) (*models.Course, error) {
    return s.Repo.GetCourseByID(id)
}

func (s *CourseService) UpdateCourse(id int, updates map[string]interface{}) (*models.Course, error) {
    return s.Repo.UpdateCourse(id, updates)
}

func (s *CourseService) DeleteCourse(id int) error {
    return s.Repo.DeleteCourse(id)
}