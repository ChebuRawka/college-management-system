package repositories

import (
	"backend/models"
	"backend/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) CreateStudent(student *models.Student) error {
	// Проверяем, существует ли курс с указанным group_name
	query := `
        SELECT id 
        FROM courses 
        WHERE name = $1
    `
	var courseID int
	err := r.DB.QueryRow(query, student.GroupName).Scan(&courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("course with name '%s' does not exist", student.GroupName)
		}
		return err
	}

	// Преобразуем дату рождения в time.Time
	dateOfBirth, err := time.Parse("2006-01-02", student.DateOfBirth)
	if err != nil {
		return fmt.Errorf("invalid date_of_birth format: %v", err)
	}

	// Вставляем данные студента в базу данных
	insertQuery := `
        INSERT INTO students (name, date_of_birth, group_name)
        VALUES ($1, $2, $3)
        RETURNING id
    `
	err = r.DB.QueryRow(insertQuery, student.Name, dateOfBirth, student.GroupName).Scan(&student.ID)
	return err
}

func (r *StudentRepository) GetStudents() ([]models.Student, error) {
    query := `
        SELECT s.id, s.name, s.date_of_birth, s.group_name, c.teacher_id
        FROM students s
        LEFT JOIN courses c ON s.group_name = c.name
    `
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var students []models.Student
    for rows.Next() {
        var student models.Student
        var dateOfBirth time.Time
        var teacherID sql.NullInt64
        if err := rows.Scan(&student.ID, &student.Name, &dateOfBirth, &student.GroupName, &teacherID); err != nil {
            return nil, err
        }

        // Преобразуем date_of_birth в строку
        student.DateOfBirth = dateOfBirth.Format("2006-01-02")

        // Вычисляем возраст
        student.Age = utils.CalculateAge(dateOfBirth)

        // Если teacher_id существует, добавляем его в ответ
        if teacherID.Valid {
            teacherIDValue := int(teacherID.Int64)
            student.TeacherID = &teacherIDValue
        } else {
            student.TeacherID = nil
        }

        students = append(students, student)
    }
    return students, nil
}

func (r *StudentRepository) GetStudentByID(id int) (*models.Student, error) {
    query := `
        SELECT s.id, s.name, s.date_of_birth, s.group_name, c.teacher_id
        FROM students s
        LEFT JOIN courses c ON s.group_name = c.name
        WHERE s.id = $1
    `
    row := r.DB.QueryRow(query, id)

    var student models.Student
    var dateOfBirth time.Time
    var teacherID sql.NullInt64
    if err := row.Scan(&student.ID, &student.Name, &dateOfBirth, &student.GroupName, &teacherID); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("student with id %d not found", id)
        }
        return nil, err
    }

    // Преобразуем date_of_birth в строку
    student.DateOfBirth = dateOfBirth.Format("2006-01-02")

    // Вычисляем возраст
    student.Age = utils.CalculateAge(dateOfBirth)

    // Если teacher_id существует, добавляем его в ответ
    if teacherID.Valid {
        teacherIDValue := int(teacherID.Int64)
        student.TeacherID = &teacherIDValue
    } else {
        student.TeacherID = nil
    }

    return &student, nil
}

func (r *StudentRepository) UpdateStudent(id int, updates map[string]interface{}) (*models.Student, error) {
	setClauses := []string{}
	args := []interface{}{}
	paramIndex := 1

	for key, value := range updates {
		switch key {
		case "name":
			setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramIndex))
			args = append(args, value)
			paramIndex++
		case "date_of_birth":
			dateOfBirth, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type for date_of_birth")
			}
			parsedDate, err := time.Parse("2006-01-02", dateOfBirth)
			if err != nil {
				return nil, fmt.Errorf("invalid date_of_birth format: %v", err)
			}
			setClauses = append(setClauses, fmt.Sprintf("date_of_birth = $%d", paramIndex))
			args = append(args, parsedDate)
			paramIndex++
		case "group_name":
			groupName, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type for group_name")
			}

			// Проверяем существование курса
			exists, err := r.CourseExists(groupName)
			if err != nil {
				return nil, err
			}
			if !exists {
				return nil, fmt.Errorf("course with name '%s' does not exist", groupName)
			}

			setClauses = append(setClauses, fmt.Sprintf("group_name = $%d", paramIndex))
			args = append(args, groupName)
			paramIndex++
		default:
			return nil, fmt.Errorf("invalid field: %s", key)
		}
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`UPDATE students SET %s WHERE id = $%d RETURNING id, name, date_of_birth, group_name`, strings.Join(setClauses, ", "), paramIndex)
	args = append(args, id)

	var student models.Student
	var dateOfBirth time.Time
	err := r.DB.QueryRow(query, args...).Scan(&student.ID, &student.Name, &dateOfBirth, &student.GroupName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student with id %d not found", id)
		}
		return nil, err
	}

	student.DateOfBirth = dateOfBirth.Format("2006-01-02")
	student.Age = utils.CalculateAge(dateOfBirth)

	return &student, nil
}

func (r *StudentRepository) DeleteStudent(id int) error {
	query := `DELETE FROM students WHERE id = $1`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("student with id %d not found", id)
	}
	return nil
}

func (r *StudentRepository) CourseExists(courseName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM courses WHERE name = $1)`
	var exists bool
	err := r.DB.QueryRow(query, courseName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
