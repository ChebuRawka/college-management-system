package repositories

import (
    "backend/models"
    "database/sql"
    "errors"
	"fmt"
	"strings"
    "time"
)

type ScheduleRepository struct {
    DB *sql.DB
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
    return &ScheduleRepository{DB: db}
}

// CreateSchedule создает новую запись в расписании
func (r *ScheduleRepository) CreateSchedule(teacherID, classroomID int, schedule *models.Schedule) error {
    query := `
        INSERT INTO schedules (teacher_id, classroom_id, group_name, start_time, end_time, day_of_week)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    err := r.DB.QueryRow(query, teacherID, classroomID, schedule.GroupName, schedule.StartTime, schedule.EndTime, schedule.DayOfWeek).Scan(&schedule.ID)
    if err != nil {
        return err
    }

    // Подтягиваем teacher_name и classroom_name для ответа
    query = `
        SELECT t.name AS teacher_name, c.name AS classroom_name
        FROM schedules s
        LEFT JOIN teachers t ON s.teacher_id = t.id
        LEFT JOIN classrooms c ON s.classroom_id = c.id
        WHERE s.id = $1
    `
    row := r.DB.QueryRow(query, schedule.ID)
    if err := row.Scan(&schedule.TeacherName, &schedule.ClassroomName); err != nil {
        return err
    }

    return nil
}

func (r *ScheduleRepository) CheckScheduleConflict(teacherID int, dayOfWeek string, startTime, endTime time.Time) (bool, error) {
    query := `
        SELECT EXISTS (
            SELECT 1
            FROM schedules
            WHERE teacher_id = $1
              AND day_of_week = $2
              AND (
                  ($3 < end_time AND $4 > start_time)
              )
        )
    `
    var exists bool
    err := r.DB.QueryRow(query, teacherID, dayOfWeek, startTime, endTime).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}

func (r *ScheduleRepository) GetSchedules() ([]models.Schedule, error) {
    query := `
        SELECT s.id, t.name AS teacher_name, c.name AS classroom_name, s.group_name, s.start_time, s.end_time, s.day_of_week
        FROM schedules s
        LEFT JOIN teachers t ON s.teacher_id = t.id
        LEFT JOIN classrooms c ON s.classroom_id = c.id
    `
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var schedules []models.Schedule
    for rows.Next() {
        var schedule models.Schedule
        if err := rows.Scan(&schedule.ID, &schedule.TeacherName, &schedule.ClassroomName, &schedule.GroupName, &schedule.StartTime, &schedule.EndTime, &schedule.DayOfWeek); err != nil {
            return nil, err
        }
        schedules = append(schedules, schedule)
    }
    return schedules, nil
}

func (r *ScheduleRepository) GetScheduleByID(id int) (*models.Schedule, error) {
    query := `
        SELECT s.id, t.name AS teacher_name, c.name AS classroom_name, s.group_name, s.start_time, s.end_time, s.day_of_week
        FROM schedules s
        LEFT JOIN teachers t ON s.teacher_id = t.id
        LEFT JOIN classrooms c ON s.classroom_id = c.id
        WHERE s.id = $1
    `
    row := r.DB.QueryRow(query, id)

    var schedule models.Schedule
    if err := row.Scan(&schedule.ID, &schedule.TeacherName, &schedule.ClassroomName, &schedule.GroupName, &schedule.StartTime, &schedule.EndTime, &schedule.DayOfWeek); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("schedule not found")
        }
        return nil, err
    }
    return &schedule, nil
}

func (r *ScheduleRepository) UpdateSchedule(id int, updates map[string]interface{}) (*models.Schedule, error) {
    setClauses := []string{}
    args := []interface{}{}
    paramIndex := 1

    // Собираем условия для обновления
    for key, value := range updates {
        switch key {
        case "teacher_id":
            setClauses = append(setClauses, fmt.Sprintf("teacher_id = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "classroom_id":
            setClauses = append(setClauses, fmt.Sprintf("classroom_id = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "group_name":
            setClauses = append(setClauses, fmt.Sprintf("group_name = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "start_time":
            setClauses = append(setClauses, fmt.Sprintf("start_time = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "end_time":
            setClauses = append(setClauses, fmt.Sprintf("end_time = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "day_of_week":
            setClauses = append(setClauses, fmt.Sprintf("day_of_week = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        default:
            return nil, errors.New("invalid field: " + key)
        }
    }

    if len(setClauses) == 0 {
        return nil, errors.New("no fields to update")
    }

    // Формируем SQL-запрос для обновления
    query := fmt.Sprintf(`UPDATE schedules SET %s WHERE id = $%d RETURNING id`, strings.Join(setClauses, ", "), paramIndex)
    args = append(args, id)

    var scheduleID int
    err := r.DB.QueryRow(query, args...).Scan(&scheduleID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("schedule not found")
        }
        return nil, err
    }

    // Подтягиваем обновленные данные с именами
    query = `
        SELECT s.id, t.name AS teacher_name, c.name AS classroom_name, s.group_name, s.start_time, s.end_time, s.day_of_week
        FROM schedules s
        LEFT JOIN teachers t ON s.teacher_id = t.id
        LEFT JOIN classrooms c ON s.classroom_id = c.id
        WHERE s.id = $1
    `
    row := r.DB.QueryRow(query, scheduleID)

    var schedule models.Schedule
    if err := row.Scan(&schedule.ID, &schedule.TeacherName, &schedule.ClassroomName, &schedule.GroupName, &schedule.StartTime, &schedule.EndTime, &schedule.DayOfWeek); err != nil {
        return nil, err
    }

    return &schedule, nil
}

// DeleteSchedule удаляет запись расписания по ID
func (r *ScheduleRepository) DeleteSchedule(id int) error {
    query := `DELETE FROM schedules WHERE id = $1`
    result, err := r.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("schedule not found")
    }
    return nil
}

func (r *ScheduleRepository) GetFilteredSchedules(dayOfWeek, groupName string) ([]models.Schedule, error) {
    query := `
        SELECT s.id, t.name AS teacher_name, c.name AS classroom_name, s.group_name, s.start_time, s.end_time, s.day_of_week
        FROM schedules s
        LEFT JOIN teachers t ON s.teacher_id = t.id
        LEFT JOIN classrooms c ON s.classroom_id = c.id
        WHERE 1=1
    `
    args := []interface{}{}
    paramIndex := 1

    if dayOfWeek != "" {
        query += fmt.Sprintf(" AND s.day_of_week = $%d", paramIndex)
        args = append(args, dayOfWeek)
        paramIndex++
    }

    if groupName != "" {
        query += fmt.Sprintf(" AND s.group_name = $%d", paramIndex)
        args = append(args, groupName)
        paramIndex++
    }

    rows, err := r.DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var schedules []models.Schedule
    for rows.Next() {
        var schedule models.Schedule
        if err := rows.Scan(&schedule.ID, &schedule.TeacherName, &schedule.ClassroomName, &schedule.GroupName, &schedule.StartTime, &schedule.EndTime, &schedule.DayOfWeek); err != nil {
            return nil, err
        }
        schedules = append(schedules, schedule)
    }
    return schedules, nil
}


