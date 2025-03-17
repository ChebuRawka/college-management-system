CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    teacher_id INT REFERENCES teachers(id) ON DELETE SET NULL
);