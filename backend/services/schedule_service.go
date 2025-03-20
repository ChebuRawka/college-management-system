package services

import (
    "backend/models"
    "backend/repository"
    "fmt"
    "errors"
)

type ScheduleService struct {
    Repo *repositories.ScheduleRepository
    TeacherRepo *repositories.TeacherRepository
}

func NewScheduleService(
    scheduleRepo *repositories.ScheduleRepository,
    teacherRepo *repositories.TeacherRepository, // Добавляем параметр для TeacherRepository
) *ScheduleService {
    return &ScheduleService{
        Repo:       scheduleRepo,
        TeacherRepo: teacherRepo,
    }
}

func (s *ScheduleService) CreateSchedule(teacherID, classroomID int, schedule *models.Schedule) error {
    // Проверка, что start_time < end_time
    if !schedule.StartTime.Before(schedule.EndTime) {
        return errors.New("start_time must be before end_time")
    }

    // Проверка, что продолжительность занятия равна 90 минутам (1.5 часа)
    duration := schedule.EndTime.Sub(schedule.StartTime)
    if duration.Minutes() != 90 {
        return errors.New("lesson duration must be exactly 1.5 hours (90 minutes)")
    }

    // Проверяем пересечение времени
    conflict, err := s.Repo.CheckScheduleConflict(teacherID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime)
    if err != nil {
        return err
    }
    if conflict {
        return errors.New("teacher already has a class at this time")
    }

    // Получаем продолжительность занятия в часах
    durationInHours := duration.Hours()
    fmt.Printf("Duration of the lesson: %.2f hours\n", durationInHours)

    // Проверяем и списываем рабочие часы у преподавателя
    err = s.TeacherRepo.UpdateTeacherWorkingHours(teacherID, durationInHours)
    if err != nil {
        fmt.Println("Error updating teacher working hours:", err)
        return err
    }

    fmt.Println("Teacher working hours updated successfully")

    // Создаем запись в расписании
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
    if groupName == "" {
        return nil, errors.New("group name cannot be empty")
    }
    schedules, err := s.Repo.GetFilteredSchedules("", groupName)
    if err != nil {
        return nil, err
    }
    return schedules, nil
}

