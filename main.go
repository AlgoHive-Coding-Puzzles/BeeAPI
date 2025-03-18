package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/algohive/beeapi/controllers"
	"github.com/algohive/beeapi/middleware"
	"github.com/algohive/beeapi/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/algohive/beeapi/docs" // Import generated docs
)

// @title BeeAPI Go
// @version 1.0
// @description API server for AlgoHive puzzles
// @contact.name API Support
// @contact.email ericphlpp@proton.me
// @license.name MIT
// @license.url https://github.com/AlgoHive-Coding-Puzzles/BeeAPI/blob/main/LICENSE
// @host localhost:5000
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Create services
	authService := services.NewAuthService()
	puzzlesLoader := services.NewPuzzlesLoader()
	pythonRunner := services.NewPythonRunner(os.Getenv("PYTHON_PATH")) // Get from env or use default

	// Load users from environment variables
	authService.LoadUsersFromEnv()

	// Create controllers
	healthController := controllers.NewHealthController()
	themeController := controllers.NewThemeController(puzzlesLoader)
	puzzleController := controllers.NewPuzzleController(puzzlesLoader, pythonRunner)
	authController := controllers.NewAuthController(authService)

	// Create router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"} // Allow all origins
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.GET("/ping", healthController.Ping)
	router.GET("/name", healthController.GetServerName)

	// Theme routes (public)
	router.GET("/themes", themeController.GetThemes)
	router.GET("/themes/names", themeController.GetThemeNames)
	router.GET("/theme", themeController.GetTheme)

	// Puzzle routes (public)
	router.GET("/puzzles", puzzleController.GetPuzzles)
	router.GET("/puzzles/names", puzzleController.GetPuzzleNames)
	router.GET("/puzzle", puzzleController.GetPuzzle)
	router.GET("/puzzle/generate", puzzleController.GeneratePuzzle)

	// Auth routes
	auth := router.Group("/auth")
	{
		// Public auth endpoints
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		
		// Protected auth endpoints
		authProtected := auth.Group("")
		authProtected.Use(middleware.RequireAuth(authService))
		{
			authProtected.POST("/logout", authController.Logout)
			authProtected.DELETE("/delete-account", authController.DeleteAccount)
			authProtected.GET("/user", authController.GetUser)
			authProtected.GET("/check", authController.CheckAuth)
		}
	}
	
	// Protected routes
	protected := router.Group("")
	protected.Use(middleware.RequireAuth(authService))
	{
		// Theme management
		protected.POST("/theme", themeController.CreateTheme)
		protected.DELETE("/theme", themeController.DeleteTheme)
		protected.POST("/theme/reload", themeController.ReloadThemes)
		
		// Puzzle management
		protected.POST("/puzzle/upload", puzzleController.UploadPuzzle)
		protected.DELETE("/puzzle", puzzleController.DeletePuzzle)
	}
	
	// // Serve static files for frontend
	// router.Static("/", "./public")
	// router.NoRoute(func(c *gin.Context) {
	// 	c.File("./public/index.html")
	// })
	
	// Make sure puzzles directory exists
	os.MkdirAll(services.PuzzlesDir, 0755)
	
	// Extract and load puzzles
	log.Println("Extracting puzzles...")
	if err := puzzlesLoader.Extract(); err != nil {
		log.Printf("Warning: Failed to extract puzzles: %v", err)
	}
	
	log.Println("Loading puzzles...")
	if err := puzzlesLoader.Load(); err != nil {
		log.Printf("Warning: Failed to load puzzles: %v", err)
	}
	
	// Determine port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	
	// Start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	
	// Run server in a goroutine so it doesn't block
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Unload puzzles
	log.Println("Unloading puzzles...")
	if err := puzzlesLoader.Unload(); err != nil {
		log.Printf("Warning: Failed to unload puzzles: %v", err)
	}
	
	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited gracefully")
}
