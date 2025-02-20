package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors" // ✅ Import CORS middleware
    "task-management-system/internal/handlers"
    "task-management-system/internal/websocket"
    "task-management-system/internal/middleware"
    "gorm.io/gorm"
    "time"
    "task-management-system/internal/config"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
    router := gin.Default()

    // ✅ Allow frontend (http://localhost:3000) to call APIs
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://your-frontend-url.vercel.app"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // ✅ WebSocket Route (Make sure WebSocket is properly set)
    router.GET("/ws", func(c *gin.Context) {
        websocket.HandleWebSocket(c.Writer, c.Request)
    })

    router.POST("/register", handlers.RegisterUser(db))
    router.POST("/login", handlers.LoginUser(db)) // ✅ Fixed reference

    // Protected routes group
    protected := router.Group("/api")
    protected.Use(middleware.JWTAuthMiddleware())
    {
        // User routes
        protected.GET("/profile", handlers.GetProfile(db))
        
        // Task routes
        protected.GET("/tasks", handlers.GetTasks(db))
        protected.POST("/tasks", handlers.CreateTask(db))
        protected.PUT("/tasks/:id/status", handlers.UpdateTaskStatus(db))
        protected.POST("/tasks/suggestions", handlers.GetTaskSuggestions(db, cfg))
        
        // Admin routes
        admin := protected.Group("")
        admin.Use(middleware.AdminMiddleware())
        {
            admin.DELETE("/tasks/:id", handlers.DeleteTask(db))
            admin.GET("/users", handlers.GetUsers(db))
        }
    }

    return router
}
