package repositories

import (
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
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
func (r *TeacherRepository) CreateTeacher(teacher *models.Teacher) error {
    // Проверяем, что все указанные курсы существуют
    if len(teacher.Courses) > 0 {
        validCourses, err := r.CheckCoursesExist(teacher.Courses)
        if err != nil {
            return fmt.Errorf("failed to validate courses: %v", err)
        }
        if !validCourses {
            return fmt.Errorf("some courses do not exist")
        }
    }

    query := `
        INSERT INTO teachers (name, subject, courses, working_hours)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    err := r.DB.QueryRow(query, teacher.Name, teacher.Subject, pq.Array(teacher.Courses), teacher.WorkingHours).Scan(&teacher.ID)
    if err != nil {
        return fmt.Errorf("failed to create teacher: %v", err)
    }
    return nil
}

// UpdateTeacherWorkingHours обновляет количество рабочих часов преподавателя
func (r *TeacherRepository) UpdateTeacherWorkingHours(teacherID int, hours float64) error {
    query := `
        UPDATE teachers
        SET working_hours = working_hours - $1
        WHERE id = $2 AND working_hours >= $1
    `
    result, err := r.DB.Exec(query, hours, teacherID)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("not enough working hours for the teacher")
    }
    return nil
}

func (r *TeacherRepository) CheckTeacherExists(name, subject string) (bool, error) {
    query := `
        SELECT EXISTS (
            SELECT 1
            FROM teachers
            WHERE name = $1 AND subject = $2
        )
    `
    var exists bool
    err := r.DB.QueryRow(query, name, subject).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}

func (r *TeacherRepository) CheckCoursesExist(courseNames []string) (bool, error) {
    query := `
        SELECT COUNT(*)
        FROM courses
        WHERE name = ANY($1)
    `
    var count int
    err := r.DB.QueryRow(query, pq.Array(courseNames)).Scan(&count)
    if err != nil {
        return false, err
    }

    return count == len(courseNames), nil
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
    query := `
        SELECT id, name, subject, courses, working_hours
        FROM teachers
    `
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var teachers []models.Teacher
    for rows.Next() {
        var teacher models.Teacher
        var courses []string
        if err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Subject, pq.Array(&courses), &teacher.WorkingHours); err != nil {
            return nil, err
        }
        teacher.Courses = courses
        teachers = append(teachers, teacher)
    }

    return teachers, nil
}

// Получение преподавателя по ID
func (r *TeacherRepository) GetTeacherByID(teacherID int) (*models.Teacher, error) {
    query := `
        SELECT id, name, subject, courses, working_hours
        FROM teachers
        WHERE id = $1
    `

    var teacher models.Teacher
    var courses []string

    err := r.DB.QueryRow(query, teacherID).Scan(
        &teacher.ID,
        &teacher.Name,
        &teacher.Subject,
        pq.Array(&courses),
        &teacher.WorkingHours,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("teacher with id %d not found", teacherID)
        }
        return nil, fmt.Errorf("failed to fetch teacher data: %v", err)
    }

    teacher.Courses = courses
    return &teacher, nil
}
func (r *TeacherRepository) UpdateTeacherPartial(id int, updates map[string]interface{}) error {
    setClauses := []string{}
    args := []interface{}{}
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
        case "working_hours":
            workingHours, ok := value.(float64) // JSON передает числа как float64
            if !ok {
                return fmt.Errorf("invalid type for working_hours: expected number")
            }
            if workingHours < 0 {
                return fmt.Errorf("working_hours cannot be negative")
            }
            setClauses = append(setClauses, fmt.Sprintf("working_hours = $%d", paramIndex))
            args = append(args, workingHours)
            paramIndex++
        default:
            return fmt.Errorf("invalid field: %s", key)
        }
    }

    if len(setClauses) == 0 {
        return fmt.Errorf("no fields to update")
    }

    query := fmt.Sprintf(`
        UPDATE teachers
        SET %s
        WHERE id = $%d
    `, strings.Join(setClauses, ", "), paramIndex)
    args = append(args, id)

    result, err := r.DB.Exec(query, args...)
    if err != nil {
        return fmt.Errorf("failed to update teacher data: %v", err)
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("teacher with id %d not found", id)
    }

    return nil
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

func (r *TeacherRepository) UpdateTeacherProfile(teacherID int, updates map[string]interface{}) error {
    if len(updates) == 0 {
        return fmt.Errorf("no fields to update")
    }

    query := "UPDATE teachers SET "
    var args []interface{}
    var setClauses []string

    for key, value := range updates {
        setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
        args = append(args, value)
    }

    // Добавляем все поля через запятую
    query += strings.Join(setClauses, ", ")

    // Добавляем условие WHERE
    query += " WHERE id = ?"
    args = append(args, teacherID)

    _, err := r.DB.Exec(query, args...)
    if err != nil {
        return err
    }

    return nil
}