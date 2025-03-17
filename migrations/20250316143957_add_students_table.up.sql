CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    group_name VARCHAR(50) NOT NULL,
    teacher_id INT REFERENCES teachers(id) ON DELETE SET NULL
);