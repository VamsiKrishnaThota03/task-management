package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"task-management-system/internal/config"
	"task-management-system/internal/handlers"
	"task-management-system/internal/middleware"
	"task-management-system/internal/websocket"
	"task-management-system/pkg/ai"
)

// SetupRouter wires all routes.
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.New()
	// /ws carries the JWT in its query string, so keep it out of access logs.
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/ws"}}), gin.Recovery())
	_ = router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false, // bearer tokens, no cookies
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	router.GET("/ws", gin.WrapF(websocket.Handler(cfg.AllowedOrigins)))

	router.POST("/register", handlers.RegisterUser(db))
	router.POST("/login", handlers.LoginUser(db))

	aiClient := ai.NewClient(cfg.GeminiKey, cfg.GeminiModel)

	protected := router.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware(db))
	{
		protected.GET("/profile", handlers.GetProfile(db))
		protected.GET("/users/assignable", handlers.GetAssignees(db))

		protected.GET("/tasks", handlers.GetTasks(db))
		protected.POST("/tasks", handlers.CreateTask(db))
		protected.POST("/tasks/suggestions", handlers.GetTaskSuggestions(aiClient))
		protected.GET("/tasks/:id", handlers.GetTask(db))
		protected.PUT("/tasks/:id", handlers.UpdateTask(db))
		protected.PUT("/tasks/:id/status", handlers.UpdateTaskStatus(db))

		admin := protected.Group("")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.DELETE("/tasks/:id", handlers.DeleteTask(db))
			admin.GET("/users", handlers.GetUsers(db))
		}
	}

	return router
}
