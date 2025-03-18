package main

import (
    "backend/config"
    "backend/handlers"
    "backend/repository"
    "backend/services"
    "backend/middleware"
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
    userRepo := repositories.NewUserRepository(db) // Добавляем репозиторий для пользователей

    // Инициализация сервиса
    teacherService := services.NewTeacherService(teacherRepo)
    studentService := services.NewStudentService(studentRepo)
    courseService := services.NewCourseService(courseRepo)
    classroomService := services.NewClassroomService(classroomRepo)
    scheduleService := services.NewScheduleService(scheduleRepo, teacherRepo) // Передаем teacherRepo
    authService := services.NewAuthService(userRepo, "your_secret_key")       // Добавляем сервис для авторизации

    // Инициализация обработчика
    teacherHandler := handlers.NewTeacherHandler(teacherService)
    studentHandler := handlers.NewStudentHandler(studentService)
    courseHandler := handlers.NewCourseHandler(courseService)
    classroomHandler := handlers.NewClassroomHandler(classroomService)
    scheduleHandler := handlers.NewScheduleHandler(scheduleService)
    authHandler := handlers.NewAuthHandler(authService) // Добавляем обработчик для авторизации

    // Роутер
    r := gin.Default()

    api := r.Group("/api")

    // Маршруты для авторизации
    api.POST("/register", authHandler.Register) // Регистрация нового пользователя
    api.POST("/login", authHandler.Login)       // Авторизация пользователя

    // Защищенные маршруты
    authorized := api.Group("/")
    authorized.Use(middleware.AuthMiddleware("your_secret_key")) // Middleware для проверки JWT-токена
    {
        // Только администраторы
        admin := authorized.Group("/")
        admin.Use(middleware.RoleMiddleware("admin"))
        {
            admin.GET("/admin", func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "welcome, admin!"})
            })

            // Маршруты только для администраторов
            admin.GET("/teachers", teacherHandler.GetAllTeachers)
            admin.POST("/teachers", teacherHandler.CreateTeacher)
            admin.PATCH("/teachers/:id", teacherHandler.UpdateTeacherPartial)
            admin.DELETE("/teachers/:id", teacherHandler.DeleteTeacher)

            admin.GET("/students", studentHandler.GetStudents)
            admin.POST("/students", studentHandler.CreateStudent)
            admin.GET("/students/:id", studentHandler.GetStudentByID)
            admin.PATCH("/students/:id", studentHandler.UpdateStudent)
            admin.DELETE("/students/:id", studentHandler.DeleteStudent)

            admin.GET("/courses", courseHandler.GetCourses)
            admin.POST("/courses", courseHandler.CreateCourse)
            admin.GET("/courses/:id", courseHandler.GetCourseByID)
            admin.PATCH("/courses/:id", courseHandler.UpdateCourse)
            admin.DELETE("/courses/:id", courseHandler.DeleteCourse)

            admin.POST("/classrooms", classroomHandler.CreateClassroom)
            admin.GET("/classrooms", classroomHandler.GetClassrooms)
            admin.GET("/classrooms/:id", classroomHandler.GetClassroomByID)
            admin.PATCH("/classrooms/:id", classroomHandler.UpdateClassroom)
            admin.DELETE("/classrooms/:id", classroomHandler.DeleteClassroom)

            admin.POST("/schedules", scheduleHandler.CreateSchedule)
            admin.GET("/schedules", scheduleHandler.GetSchedules)
            admin.GET("/schedules/:id", scheduleHandler.GetScheduleByID)
            admin.PATCH("/schedules/:id", scheduleHandler.UpdateSchedule)
            admin.DELETE("/schedules/:id", scheduleHandler.DeleteSchedule)
        }

        // Учителя и администраторы
    teacher := authorized.Group("/")
    teacher.Use(middleware.RoleMiddleware("teacher", "admin"))
    {
        // Маршрут для просмотра расписания учителя
        teacher.GET("/teachers/:teacher_name/schedule", teacherHandler.GetTeacherSchedule)
    }
    }

    r.Run(":8080")
}