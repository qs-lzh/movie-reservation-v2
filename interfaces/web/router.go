package web

import (
	"github.com/gin-gonic/gin"

	"github.com/qs-lzh/movie-reservation/internal/app"
	"github.com/qs-lzh/movie-reservation/internal/handler"
	"github.com/qs-lzh/movie-reservation/internal/middleware"
)

func InitRouter(app *app.App) *gin.Engine {
	authHandler := handler.NewAuthHandler(app)
	movieHandler := handler.NewMovieHandler(app)
	showtimeHandler := handler.NewShowtimeHandler(app)
	reservationHandler := handler.NewReservationHandler(app)
	hallHandler := handler.NewHallHandler(app)
	captchaHandler := handler.NewCaptchaHandler(app)

	r := gin.New()

	// use zap as logger
	r.Use(
		middleware.ZapLogger(app.Logger),
		middleware.ErrorLogger(app.Logger),
	)

	r.GET("/captcha", captchaHandler.GenerateCaptcha)
	r.POST("/captcha", captchaHandler.VerifyCaptcha)

	users := r.Group("/users")
	{
		// [User] [Admin]
		users.POST("/register", authHandler.Register)
		users.POST("/login", authHandler.Login)
		users.POST("/logout", authHandler.Logout)
	}

	movies := r.Group("movies")
	{
		movies.GET("/", movieHandler.GetAllMovies)
		movies.GET("/:id", movieHandler.GetMovieByID)
		movies.GET("/:id/showtimes", movieHandler.GetMovieShowtimes)
		// [Admin]
		movies.POST("/", middleware.RequireAuth(), middleware.RequireAdmin(), movieHandler.CreateMovie)
		movies.PUT("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), movieHandler.UpdateMovie)
		movies.DELETE("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), movieHandler.DeleteMovie)
	}

	showtimes := r.Group("showtimes")
	{
		showtimes.GET("/", showtimeHandler.ListAllShowtimes)
		showtimes.GET("/:id", showtimeHandler.GetShowtimeByID)
		showtimes.GET("/:id/availability", showtimeHandler.GetShowtimeAvailability)
		// [Admin]
		showtimes.POST("/", middleware.RequireAuth(), middleware.RequireAdmin(), showtimeHandler.CreateShowtime)
		showtimes.PUT("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), showtimeHandler.UpdateShowtime)
		showtimes.DELETE("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), showtimeHandler.DeleteShowtimeByID)
	}

	reservations := r.Group("reservations")
	reservations.Use(middleware.RequireAuth())
	{
		// [User]
		reservations.POST("/", reservationHandler.CreateReservation)
		reservations.GET("/me", reservationHandler.GetMyReservations)
		reservations.DELETE("/:id", reservationHandler.CancelReservation)
	}

	halls := r.Group("halls")
	{
		halls.GET("/", hallHandler.GetAllHalls)
		halls.GET("/:id", hallHandler.GetHallByID)
		// [Admin]
		halls.POST("/", middleware.RequireAuth(), middleware.RequireAdmin(), hallHandler.CreateHall)
		halls.PUT("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), hallHandler.UpdateHall)
		halls.DELETE("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), hallHandler.DeleteHall)
	}

	return r
}
