package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rimvydascivilis/book-tracker/backend/config"
	mariadbRepo "github.com/rimvydascivilis/book-tracker/backend/internal/repository/mariadb"
	"github.com/rimvydascivilis/book-tracker/backend/internal/rest"
	"github.com/rimvydascivilis/book-tracker/backend/services/auth"
	"github.com/rimvydascivilis/book-tracker/backend/services/book"
	"github.com/rimvydascivilis/book-tracker/backend/services/goal"
	"github.com/rimvydascivilis/book-tracker/backend/services/list"
	"github.com/rimvydascivilis/book-tracker/backend/services/note"
	"github.com/rimvydascivilis/book-tracker/backend/services/progress"
	"github.com/rimvydascivilis/book-tracker/backend/services/reading"
	"github.com/rimvydascivilis/book-tracker/backend/services/user"
	"github.com/rimvydascivilis/book-tracker/backend/services/validation"
	"github.com/rimvydascivilis/book-tracker/backend/utils"
)

func main() {
	cfg := config.LoadConfig()

	utils.SetupLogger(cfg.LogLevel)

	// Database connection
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Europe/Vilnius")
	dsn := fmt.Sprintf("%s?%s", cfg.DBUrl, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		utils.Fatal("failed to open connection to database", err)
	}
	err = dbConn.Ping()
	if err != nil {
		utils.Fatal("failed to ping database", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			utils.Fatal("got error when closing the DB connection", err)
		}
	}()

	e := echo.New()

	// Middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			msg := fmt.Sprintf("request timed out on route %s", c.Path())
			utils.Error(msg, err)
		},
		Timeout: 5 * time.Second,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			utils.Info("request completed", map[string]interface{}{
				"method": c.Request().Method,
				"uri":    values.URI,
				"status": values.Status,
			})

			return nil
		},
	}))

	// Repositories
	userRepo := mariadbRepo.NewUserRepository(dbConn)
	bookRepo := mariadbRepo.NewBookRepository(dbConn)
	goalRepo := mariadbRepo.NewGoalRepository(dbConn)
	readingRepo := mariadbRepo.NewReadingRepository(dbConn)
	progressRepo := mariadbRepo.NewProgressRepository(dbConn)
	listRepo := mariadbRepo.NewListRepository(dbConn)
	listItemRepo := mariadbRepo.NewListItemRepository(dbConn)
	noteRepo := mariadbRepo.NewNoteRepository(dbConn)

	// Services
	validationSvc := validation.NewValidationService()
	googleOauth2Svc, err := auth.NewGoogleOAuth2Service()
	if err != nil {
		utils.Fatal("failed to create Google OAuth2 service", err)
	}
	jwtSvc := auth.NewJWTService(cfg.JWTSecret, userRepo)
	userSvc := user.NewUserService(userRepo, validationSvc)
	authSvc := auth.NewAuthService(userSvc, googleOauth2Svc, jwtSvc)
	bookSvc := book.NewBookService(bookRepo, validationSvc)
	goalSvc := goal.NewGoalService(goalRepo, progressRepo, readingRepo, validationSvc)
	readingSvc := reading.NewReadingService(readingRepo, progressRepo, bookRepo, validationSvc)
	progressSvc := progress.NewProgressService(progressRepo, readingRepo, validationSvc)
	listSvc := list.NewListService(listRepo, listItemRepo, bookRepo, validationSvc)
	noteSvc := note.NewNoteService(bookRepo, noteRepo, validationSvc)

	// Handlers
	authH := rest.NewAuthHandler(authSvc)
	bookH := rest.NewBookHandler(bookSvc)
	goalH := rest.NewGoalHandler(goalSvc)
	readingH := rest.NewReadingHandler(readingSvc)
	progressH := rest.NewProgressHandler(progressSvc)
	listH := rest.NewListHandler(listSvc)
	noteH := rest.NewNoteHandler(noteSvc)

	// Route groups
	api := e.Group("/api")
	authenticatedApi := api.Group("", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
		ContextKey: "user",
	}))

	// Unauthenticated routes
	api.POST("/auth/login", authH.Login)

	// Authenticated routes
	authenticatedApi.GET("/books", bookH.GetBooks)
	authenticatedApi.GET("/books/search", bookH.SearchBooks) // ?title=My%20book
	authenticatedApi.POST("/books", bookH.CreateBook)
	authenticatedApi.PUT("/books/:id", bookH.UpdateBook)
	authenticatedApi.DELETE("/books/:id", bookH.DeleteBook)

	authenticatedApi.GET("/goal", goalH.GetGoal)
	authenticatedApi.GET("/goal/progress", goalH.GetGoalProgress)
	authenticatedApi.PUT("/goal", goalH.SetGoal)

	authenticatedApi.GET("/readings", readingH.GetReadings)
	authenticatedApi.POST("/readings", readingH.CreateReading)

	authenticatedApi.POST("/progress/:readingId", progressH.CreateProgress)

	authenticatedApi.GET("/lists", listH.ListLists)
	authenticatedApi.GET("/list", listH.GetList)                                      // ?list_id=1
	authenticatedApi.POST("/list", listH.CreateList)                                  // {"title": "My list"}
	authenticatedApi.POST("/list/item", listH.AddBookToList)                          // {"list_id": 1, "book_id": 1}
	authenticatedApi.DELETE("/list/:list_id/item/:item_id", listH.RemoveBookFromList) // /list/1/item/1

	authenticatedApi.GET("/notes/:book_id", noteH.GetNotes)      // /notes/1
	authenticatedApi.POST("/notes/:book_id", noteH.CreateNote)   // /notes/1 {"page_number": 1, "content": "My note"}
	authenticatedApi.DELETE("/notes/:note_id", noteH.DeleteNote) // /notes/1

	log.Fatal(e.Start(cfg.ServerAddr))
}
