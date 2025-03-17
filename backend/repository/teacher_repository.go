package repositories

import (
    "backend/models"
    "database/sql"
	"strings"
    "fmt"
    "errors"
)

type TeacherRepository struct {
    DB *sql.DB
    CourseRepo *CourseRepository
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
    return &TeacherRepository{
        DB:         db,
        CourseRepo: NewCourseRepository(db), // Инициализируем CourseRepository
    }
}

// Создание преподавателя
func (r *TeacherRepository) CreateTeacher(teacher models.Teacher) (int, error) {
    query := `INSERT INTO teachers (name, subject) VALUES ($1, $2) RETURNING id`
    var id int
    err := r.DB.QueryRow(query, teacher.Name, teacher.Subject).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (r *TeacherRepository) GetAllTeachersWithCourses() ([]models.Teacher, error) {
    query := `
        SELECT t.id, t.name, t.subject, c.name AS course_name
        FROM teachers t
        LEFT JOIN courses c ON t.id = c.teacher_id
    `
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var teachers []models.Teacher
    teacherMap := make(map[int]*models.Teacher)

    for rows.Next() {
        var teacherID int
        var teacherName, teacherSubject, courseName sql.NullString
        if err := rows.Scan(&teacherID, &teacherName, &teacherSubject, &courseName); err != nil {
            return nil, err
        }

        // Если преподаватель ещё не добавлен в map, добавляем его
        if _, exists := teacherMap[teacherID]; !exists {
            teacherMap[teacherID] = &models.Teacher{
                ID:       teacherID,
                Name:     teacherName.String,
                Subject:  teacherSubject.String,
                Courses:  []string{},
            }
        }

        // Если у преподавателя есть группа, добавляем её название в список
        if courseName.Valid {
            teacherMap[teacherID].Courses = append(teacherMap[teacherID].Courses, courseName.String)
        }
    }

    // Преобразуем map в slice
    for _, teacher := range teacherMap {
        teachers = append(teachers, *teacher)
    }

    return teachers, nil
}
// Получение всех преподавателей
func (r *TeacherRepository) GetAllTeachers() ([]models.Teacher, error) {
    query := `SELECT id, name, subject FROM teachers`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var teachers []models.Teacher
    for rows.Next() {
        var teacher models.Teacher
        err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Subject)
        if err != nil {
            return nil, err
        }
        teachers = append(teachers, teacher)
    }
    return teachers, nil
}

// Получение преподавателя по ID
func (r *TeacherRepository) GetTeacherByID(id int) (*models.Teacher, error) {
    query := `SELECT id, name, subject FROM teachers WHERE id = $1`
    row := r.DB.QueryRow(query, id)

    var teacher models.Teacher
    err := row.Scan(&teacher.ID, &teacher.Name, &teacher.Subject)
    if err == sql.ErrNoRows {
        return nil, nil // Преподаватель не найден
    }
    if err != nil {
        return nil, err
    }
    return &teacher, nil
}

// Частичное обновление преподавателя (PATCH)
func (r *TeacherRepository) UpdateTeacherPartial(id int, updates map[string]interface{}) (*models.Teacher, error) {
    // Проверяем, существует ли запись
    exists, err := r.TeacherExists(id)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, fmt.Errorf("teacher with id %d not found", id)
    }

    // Формируем SQL-запрос
    query := `UPDATE teachers SET `
    args := []interface{}{}
    setClauses := []string{}
    paramIndex := 1

    for key, value := range updates {
        switch key {
        case "name":
            setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "subject":
            setClauses = append(setClauses, fmt.Sprintf("subject = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        default:
            return nil, fmt.Errorf("invalid field: %s", key)
        }
    }

    if len(setClauses) == 0 {
        return nil, fmt.Errorf("no fields to update")
    }

    query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE id = $%d RETURNING *", paramIndex)
    args = append(args, id)

    // Выполняем запрос и получаем обновлённую запись
    var teacher models.Teacher
    err = r.DB.QueryRow(query, args...).Scan(&teacher.ID, &teacher.Name, &teacher.Subject)
    if err != nil {
        return nil, err
    }

    return &teacher, nil
}
func (r *TeacherRepository) TeacherExists(id int) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM teachers WHERE id = $1)`
    err := r.DB.QueryRow(query, id).Scan(&exists)
    return exists, err
}

func (r *TeacherRepository) DeleteTeacher(id int) error {
    // Проверяем, существует ли запись
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM teachers WHERE id = $1)`
    err := r.DB.QueryRow(query, id).Scan(&exists)
    if err != nil {
        return err
    }

    if !exists {
        return fmt.Errorf("teacher with id %d not found", id)
    }

    // Удаляем запись
    query = `DELETE FROM teachers WHERE id = $1`
    _, err = r.DB.Exec(query, id)
    return err
}

func (r *TeacherRepository) GetTeacherSchedule(teacherName string) ([]models.ScheduleResponse, error) {
    query := `
        SELECT s.id, t.name AS teacher_name, c.name AS classroom_name, s.group_name, s.start_time, s.end_time, s.day_of_week
        FROM schedules s
        LEFT JOIN teachers t ON s.teacher_id = t.id
        LEFT JOIN classrooms c ON s.classroom_id = c.id
        WHERE t.name = $1
    `
    rows, err := r.DB.Query(query, teacherName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var schedules []models.ScheduleResponse
    for rows.Next() {
        var schedule models.ScheduleResponse
        if err := rows.Scan(&schedule.ID, &schedule.TeacherName, &schedule.ClassroomName, &schedule.GroupName, &schedule.StartTime, &schedule.EndTime, &schedule.DayOfWeek); err != nil {
            return nil, err
        }
        schedules = append(schedules, schedule)
    }

    if len(schedules) == 0 {
        return nil, errors.New("teacher not found")
    }
    return schedules, nil
}