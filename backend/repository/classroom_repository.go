package repositories

import (
    "backend/models"
    "database/sql"
    "errors"
	"fmt"
	"strings"
)

type ClassroomRepository struct {
    DB *sql.DB
}

func NewClassroomRepository(db *sql.DB) *ClassroomRepository {
    return &ClassroomRepository{DB: db}
}

// CreateClassroom создает новую аудиторию
func (r *ClassroomRepository) CreateClassroom(classroom *models.Classroom) error {
    query := `
        INSERT INTO classrooms (name, capacity, description)
        VALUES ($1, $2, $3)
        RETURNING id
    `
    err := r.DB.QueryRow(query, classroom.Name, classroom.Capacity, classroom.Description).Scan(&classroom.ID)
    return err
}

// GetClassrooms возвращает все аудитории
func (r *ClassroomRepository) GetClassrooms() ([]models.Classroom, error) {
    query := `SELECT id, name, capacity, description FROM classrooms`
    rows, err := r.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var classrooms []models.Classroom
    for rows.Next() {
        var classroom models.Classroom
        if err := rows.Scan(&classroom.ID, &classroom.Name, &classroom.Capacity, &classroom.Description); err != nil {
            return nil, err
        }
        classrooms = append(classrooms, classroom)
    }
    return classrooms, nil
}

// GetClassroomByID возвращает аудиторию по ID
func (r *ClassroomRepository) GetClassroomByID(id int) (*models.Classroom, error) {
    query := `SELECT id, name, capacity, description FROM classrooms WHERE id = $1`
    row := r.DB.QueryRow(query, id)

    var classroom models.Classroom
    if err := row.Scan(&classroom.ID, &classroom.Name, &classroom.Capacity, &classroom.Description); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("classroom not found")
        }
        return nil, err
    }
    return &classroom, nil
}

// UpdateClassroom обновляет данные аудитории
func (r *ClassroomRepository) UpdateClassroom(id int, updates map[string]interface{}) (*models.Classroom, error) {
    setClauses := []string{}
    args := []interface{}{}
    paramIndex := 1

    for key, value := range updates {
        switch key {
        case "name":
            setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramIndex))
            args = append(args, value)
            paramIndex++
        case "capacity":
            capacity, ok := value.(float64) // JSON передает числа как float64
            if !ok {
                return nil, errors.New("invalid type for capacity")
            }
            setClauses = append(setClauses, fmt.Sprintf("capacity = $%d", paramIndex))
            args = append(args, int(capacity)) // Преобразуем float64 в int
            paramIndex++
        case "description":
            description, ok := value.(string)
            if !ok {
                return nil, errors.New("invalid type for description")
            }
            setClauses = append(setClauses, fmt.Sprintf("description = $%d", paramIndex))
            args = append(args, description)
            paramIndex++
        default:
            return nil, errors.New("invalid field: " + key)
        }
    }

    if len(setClauses) == 0 {
        return nil, errors.New("no fields to update")
    }

    query := fmt.Sprintf(`UPDATE classrooms SET %s WHERE id = $%d RETURNING id, name, capacity, description`, strings.Join(setClauses, ", "), paramIndex)
    args = append(args, id)

    var classroom models.Classroom
    err := r.DB.QueryRow(query, args...).Scan(&classroom.ID, &classroom.Name, &classroom.Capacity, &classroom.Description)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("classroom not found")
        }
        return nil, err
    }
    return &classroom, nil
}

// DeleteClassroom удаляет аудиторию по ID
func (r *ClassroomRepository) DeleteClassroom(id int) error {
    query := `DELETE FROM classrooms WHERE id = $1`
    result, err := r.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("classroom not found")
    }
    return nil
}