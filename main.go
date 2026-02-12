package main

import (
	"car-management-system/config"
	"car-management-system/handlers"
	"car-management-system/middleware"
	"car-management-system/repositories"
	"context"
	"os" // <--- ADDED THIS IMPORT

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	r := gin.Default()

	// CORS Configuration
	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"*"},
	}
	r.Use(cors.New(corsConfig))

	// Load Config & Database
	err := loadConfig()
	if err != nil {
		panic(err)
	}
	connection, err := connectToDb()
	if err != nil {
		panic(err)
	}

	// ==========================================
	// INIT REPOSITORIES & HANDLERS
	// ==========================================
	authRepository := repositories.NewAuthRepository(connection)
	passwordResetRepository := repositories.NewPasswordResetRepository(connection)

	// Handlers
	authHandler := handlers.NewAuthHandler(authRepository)
	passwordResetHandler := handlers.NewPasswordResetHandler(authRepository, passwordResetRepository)

	carRepository := repositories.NewCarRepository(connection)
	carHandler := handlers.NewCarHandlers(carRepository)

	saleRepository := repositories.NewSaleRepository(connection)
	saleHandler := handlers.NewSaleHandler(carRepository, saleRepository)

	rentalRepository := repositories.NewRentalRepository(connection)
	rentalHandler := handlers.NewRentalHandler(carRepository, rentalRepository)

	paymentRepository := repositories.NewPaymentRepository(connection)
	paymentHandler := handlers.NewPaymentHandler(paymentRepository)

	maintenanceRepository := repositories.NewMaintenanceRepository(connection)
	maintenanceHandler := handlers.NewMaintenanceHandler(carRepository, maintenanceRepository)

	fuelRepository := repositories.NewFuelRepository(connection)
	fuelHandler := handlers.NewFuelHandler(carRepository, fuelRepository)

	expenseRepository := repositories.NewExpenseRepository(connection)
	expenseHandler := handlers.NewExpenseHandler(carRepository, expenseRepository)

	documentRepository := repositories.NewDocumentRepository(connection)
	documentHandler := handlers.NewDocumentHandler(carRepository, documentRepository)

	reviewRepository := repositories.NewReviewRepository(connection)
	reviewHandler := handlers.NewReviewHandler(carRepository, reviewRepository)

	adminRepository := repositories.NewAdminRepository(connection)
	adminHandler := handlers.NewAdminHandler(adminRepository)

	// AI Recommendation Handler
	recommendationHandler := handlers.NewRecommendationHandler(carRepository)

	// Serve static files
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")

	// ==========================================
	// 1. PUBLIC ROUTES (No Token Needed)
	// ==========================================
	r.POST("/api/user/signUp", authHandler.SignUp)
	r.POST("/api/user/signIn", authHandler.SignIn)
	r.POST("/api/user/signOut", authHandler.SignOut)

	r.POST("/api/user/forgot-password", passwordResetHandler.ForgotPassword)
	r.POST("/api/user/reset-password", passwordResetHandler.ResetPassword)

	// AI Recommendation Routes (Public)
	r.GET("/recommendation/questions", recommendationHandler.GetQuestions)
	r.POST("/recommendation/result", recommendationHandler.GetRecommendation)

	// ==========================================
	// 2. AUTHORIZED ROUTES (Token Required)
	// ==========================================
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	// Admin Routes
	admin := authorized.Group("/api/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/users", authHandler.GetAll)
		admin.GET("/payments", paymentHandler.GetAllPayments)
		admin.GET("/cars", carHandler.GetAllCars)
		admin.GET("/rentals", rentalHandler.GetAllRentals)
		admin.GET("/sales", saleHandler.GetAllSales)
		admin.GET("/stats", adminHandler.GetStats)
	}

	// User Routes
	{
		authorized.GET("/api/user/balance", authHandler.GetUserBalance)

		authorized.POST("/api/cars", carHandler.CreateCar)
		authorized.POST("/api/car/buy", carHandler.BuyCar)
		authorized.GET("/api/cars", carHandler.GetAllCars)
		authorized.PUT("/api/cars/:id", carHandler.UpdateCar)
		authorized.PATCH("/api/cars/:id", carHandler.PatchCar)
		authorized.DELETE("/api/cars/:id", carHandler.DeleteCar)

		authorized.POST("/api/sales", saleHandler.CreateSale)

		authorized.POST("/api/user/rentals", rentalHandler.CreateRental)
		authorized.GET("/api/user/rentals", rentalHandler.GetAllRentals)
		authorized.PUT("/api/rentals/:id/complete", rentalHandler.CompleteRental)
		authorized.PUT("/api/rentals/:id/cancel", rentalHandler.CancelRental)

		authorized.POST("/api/payments", paymentHandler.CreatePayment)
		authorized.GET("/api/payments", paymentHandler.GetMyPayments)

		authorized.POST("/api/cars/:id/maintenance", maintenanceHandler.CreateMaintenance)
		authorized.PATCH("api/maintenance/:id", maintenanceHandler.PatchMaintenance)
		authorized.GET("/api/cars/:id/maintenance", maintenanceHandler.GetMaintenance)

		authorized.POST("/api/cars/:id/fuel", fuelHandler.CreateFuel)
		authorized.PATCH("api/fuel/:id", fuelHandler.PatchFuel)
		authorized.GET("/api/cars/:id/fuel", fuelHandler.GetFuel)

		authorized.POST("/api/cars/:id/expenses", expenseHandler.CreateExpense)
		authorized.PATCH("/api/expenses/:id", expenseHandler.PatchExpense)
		authorized.GET("/api/cars/:id/expenses", expenseHandler.GetExpense)

		authorized.POST("/api/cars/:id/documents", documentHandler.CreateDocument)
		authorized.PATCH("/api/documents/:id", documentHandler.PatchDocument)
		authorized.GET("/api/cars/:id/documents", documentHandler.GetByCar)
		authorized.DELETE("/api/documents/:id", documentHandler.DeleteDocument)

		authorized.POST("/api/cars/:id/reviews", reviewHandler.CreateReview)
		authorized.GET("/api/cars/:id/reviews", reviewHandler.GetByCar)
	}

	r.Run(config.Config.AppHost)
}

func connectToDb() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func loadConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	// Set Defaults
	if err := viper.BindEnv("APP_HOST"); err != nil {
		viper.SetDefault("APP_HOST", ":8000")
	}
	if err := viper.BindEnv("DB_CONNECTION_STRING"); err != nil {
		viper.SetDefault("DB_CONNECTION_STRING", "postgres://car-management-system:1234@localhost:5414/car-management-system?sslmode=disable")
	}
	if err := viper.BindEnv("JWT_SECRET_KEY"); err != nil {
		viper.SetDefault("JWT_SECRET_KEY", "secretkey")
	}
	if err := viper.BindEnv("JWT_EXPIRES_IN"); err != nil {
		viper.SetDefault("JWT_EXPIRES_IN", "24h")
	}

	// Read the .env file
	err := viper.ReadInConfig()
	if err != nil {
		// It is okay if .env doesn't exist (e.g. production), but we return error if format is bad
		// If you want to allow missing .env, verify err type. For now, we return err.
		return err
	}

	// ⬇️ CRITICAL UPDATE: Bridge Viper config to OS Environment for the AI Handler
	// Viper reads the file, but the Google SDK looks for 'os.Getenv'
	if apiKey := viper.GetString("GEMINI_API_KEY"); apiKey != "" {
		os.Setenv("GEMINI_API_KEY", apiKey)
	}

	var mapConfig config.MapConfig
	err = viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
	}
	config.Config = &mapConfig
	return nil
}
