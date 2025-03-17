package main

import (
    "backend/config"
    "backend/handlers"
    "backend/repository"
    "backend/services"
    "github.com/gin-gonic/gin"
)

func main() {
    // Подключение к бдхе
    db := config.ConnectDB()
    defer db.Close()

    // Инициализация репозитория
    teacherRepo := repositories.NewTeacherRepository(db)
    studentRepo := repositories.NewStudentRepository(db)
    courseRepo := repositories.NewCourseRepository(db)
    classroomRepo := repositories.NewClassroomRepository(db)
    scheduleRepo := repositories.NewScheduleRepository(db)

    // Инициализация сервиса
    teacherService := services.NewTeacherService(teacherRepo)
    studentService := services.NewStudentService(studentRepo)
    courseService := services.NewCourseService(courseRepo)
    classroomService := services.NewClassroomService(classroomRepo)
    scheduleService := services.NewScheduleService(scheduleRepo)


    // Инициализация обработчика
    teacherHandler := handlers.NewTeacherHandler(teacherService)
    studentHandler := handlers.NewStudentHandler(studentService)
    courseHandler := handlers.NewCourseHandler(courseService)
    classroomHandler := handlers.NewClassroomHandler(classroomService)
    scheduleHandler := handlers.NewScheduleHandler(scheduleService)


    // Роутер
    r := gin.Default()

    api := r.Group("/api")

    // Маршруты
    api.GET("/teachers", teacherHandler.GetAllTeachers)
    api.POST("/teachers", teacherHandler.CreateTeacher)
    api.PATCH("/teachers/:id", teacherHandler.UpdateTeacherPartial)
    api.DELETE("/teachers/:id", teacherHandler.DeleteTeacher)
    api.GET("/teachers/:teacher_name/schedule", teacherHandler.GetTeacherSchedule)

    api.GET("/students", studentHandler.GetStudents)
    api.POST("/students", studentHandler.CreateStudent)
    api.GET("/students/:id", studentHandler.GetStudentByID)
    api.PATCH("/students/:id", studentHandler.UpdateStudent)
    api.DELETE("/students/:id", studentHandler.DeleteStudent)

    api.GET("/courses", courseHandler.GetCourses)
    api.POST("/courses", courseHandler.CreateCourse)
    api.GET("/courses/:id", courseHandler.GetCourseByID)
    api.PATCH("/courses/:id", courseHandler.UpdateCourse)
    api.DELETE("/courses/:id", courseHandler.DeleteCourse)

    api.POST("/classrooms", classroomHandler.CreateClassroom)
    api.GET("/classrooms", classroomHandler.GetClassrooms)
    api.GET("/classrooms/:id", classroomHandler.GetClassroomByID)
    api.PATCH("/classrooms/:id", classroomHandler.UpdateClassroom)
    api.DELETE("/classrooms/:id", classroomHandler.DeleteClassroom)

    api.POST("/schedules", scheduleHandler.CreateSchedule)
    api.GET("/schedules", scheduleHandler.GetSchedules)
    api.GET("/schedules/:id", scheduleHandler.GetScheduleByID)
    api.PATCH("/schedules/:id", scheduleHandler.UpdateSchedule)
    api.DELETE("/schedules/:id", scheduleHandler.DeleteSchedule)
    api.GET("/schedules/day/:day", scheduleHandler.GetSchedulesByDay)        // По дню недели
    api.GET("/schedules/group/:group", scheduleHandler.GetSchedulesByGroup) // По группе
        
    r.Run(":8080")
}