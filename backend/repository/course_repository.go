package repositories

import (
    "database/sql"
    "errors"
    "fmt"
    "backend/models"
	"strings"
)

type CourseRepository struct {
    DB *sql.DB
}



func NewCourseRepository(db *sql.DB) *CourseRepository {
    return &CourseRepository{DB: db}
}



// CreateCourse создаёт новый курс
func (r *CourseRepository) CreateCourse(course *models.Course) error {
    query := `INSERT INTO courses (name, description, teacher_id) VALUES ($1, $2, $3) RETURNING id`
    err := r.DB.QueryRow(query, course.Name, course.Description, course.TeacherID).Scan(&course.ID)
    return err
}



// GetCourses возвращает все курсы
func (r *CourseRepository) GetCourses() ([]models.Course, error) {
    query := `SELECT id, name, description, teacher_id FROM courses`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var courses []models.Course
    for rows.Next() {
        var course models.Course
        var teacherID sql.NullInt64
        if err := rows.Scan(&course.ID, &course.Name, &course.Description, &teacherID); err != nil {
            return nil, err
        }
        if teacherID.Valid {
            teacherIDValue := int(teacherID.Int64)
            course.TeacherID = &teacherIDValue
        } else {
            course.TeacherID = nil
        }
        courses = append(courses, course)
    }
    return courses, nil
}

// GetCourseByID возвращает курс по ID
func (r *CourseRepository) GetCourseByID(id int) (*models.Course, error) {
    query := `SELECT id, name, description, teacher_id FROM courses WHERE id = $1`
    row := r.DB.QueryRow(query, id)

    var course models.Course
    var teacherID sql.NullInt64
    if err := row.Scan(&course.ID, &course.Name, &course.Description, &teacherID); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("course with id %d not found", id)
        }
        return nil, err
    }
    if teacherID.Valid {
        teacherIDValue := int(teacherID.Int64)
        course.TeacherID = &teacherIDValue
    } else {
        course.TeacherID = nil
    }
    return &course, nil
}

// UpdateCourse обновляет данные курса
func (r *CourseRepository) UpdateCourse(id int, updates map[string]interface{}) (*models.Course, error) {
    setClauses := []string{}
    args := []interface{}{}
    paramIndex := 1

    for key, value := range updates {
        switch key {
        case "name":
            setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "description":
            setClauses = append(setClauses, fmt.Sprintf("description = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "teacher_id":
            setClauses = append(setClauses, fmt.Sprintf("teacher_id = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        default:
            return nil, fmt.Errorf("invalid field: %s", key)
        }
    }

    if len(setClauses) == 0 {
        return nil, fmt.Errorf("no fields to update")
    }

    query := fmt.Sprintf(`UPDATE courses SET %s WHERE id = $%d RETURNING id, name, description, teacher_id`, strings.Join(setClauses, ", "), paramIndex)
    args = append(args, id)

    var course models.Course
    var teacherID sql.NullInt64
    err := r.DB.QueryRow(query, args...).Scan(&course.ID, &course.Name, &course.Description, &teacherID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("course with id %d not found", id)
        }
        return nil, err
    }
    if teacherID.Valid {
        teacherIDValue := int(teacherID.Int64)
        course.TeacherID = &teacherIDValue
    } else {
        course.TeacherID = nil
    }
    return &course, nil
}

// DeleteCourse удаляет курс по ID
func (r *CourseRepository) DeleteCourse(id int) error {
    query := `DELETE FROM courses WHERE id = $1`
    result, err := r.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("course with id %d not found", id)
    }
    return nil
}

// GetCoursesByTeacherID возвращает курсы, связанные с преподавателем
func (r *CourseRepository) GetCoursesByTeacherID(teacherID int) ([]models.Course, error) {
    query := `SELECT id, name, description, teacher_id FROM courses WHERE teacher_id = $1`
    rows, err := r.DB.Query(query, teacherID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var courses []models.Course
    for rows.Next() {
        var course models.Course
        var teacherID sql.NullInt64
        if err := rows.Scan(&course.ID, &course.Name, &course.Description, &teacherID); err != nil {
            return nil, err
        }
        if teacherID.Valid {
            teacherIDValue := int(teacherID.Int64)
            course.TeacherID = &teacherIDValue
        } else {
            course.TeacherID = nil
        }
        courses = append(courses, course)
    }
    return courses, nil
}

