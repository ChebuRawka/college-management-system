package services

import (
    "backend/models"
    "backend/repository"
)

type ScheduleService struct {
    Repo *repositories.ScheduleRepository
}

func NewScheduleService(repo *repositories.ScheduleRepository) *ScheduleService {
    return &ScheduleService{Repo: repo}
}

func (s *ScheduleService) CreateSchedule(teacherID, classroomID int, schedule *models.Schedule) error {
    return s.Repo.CreateSchedule(teacherID, classroomID, schedule)
}

func (s *ScheduleService) GetSchedules() ([]models.Schedule, error) {
    return s.Repo.GetSchedules()
}

func (s *ScheduleService) GetScheduleByID(id int) (*models.Schedule, error) {
    return s.Repo.GetScheduleByID(id)
}
func (s *ScheduleService) UpdateSchedule(id int, updates map[string]interface{}) (*models.Schedule, error) {
    return s.Repo.UpdateSchedule(id, updates)
}

func (s *ScheduleService) DeleteSchedule(id int) error {
    return s.Repo.DeleteSchedule(id)
}

func (s *ScheduleService) GetFilteredSchedules(dayOfWeek, groupName string) ([]models.Schedule, error) {
    return s.Repo.GetFilteredSchedules(dayOfWeek, groupName)
}

func (s *ScheduleService) GetSchedulesByDay(dayOfWeek string) ([]models.Schedule, error) {
    return s.Repo.GetFilteredSchedules(dayOfWeek, "")
}

func (s *ScheduleService) GetSchedulesByGroup(groupName string) ([]models.Schedule, error) {
    return s.Repo.GetFilteredSchedules("", groupName)
}

