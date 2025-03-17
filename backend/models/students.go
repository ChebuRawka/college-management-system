package models


type Student struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
	DateOfBirth string    `json:"date_of_birth"`
    Age       int    `json:"age"`
    GroupName string `json:"group_name"` 
    TeacherID   *int    `json:"teacher_id"`
}