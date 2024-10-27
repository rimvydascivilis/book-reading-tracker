package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"book-tracker/config"
	mariadbRepo "book-tracker/internal/repository/mariadb"
	"book-tracker/internal/rest"
	"book-tracker/services/auth"
	"book-tracker/utils"
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
				"uri":    values.URI,
				"status": values.Status,
			})

			return nil
		},
	}))

	// Repositories
	userRepo := mariadbRepo.NewUserRepository(dbConn)

	// Services
	authSvc, err := auth.NewAuthService(userRepo, cfg.JWTSecret)
	if err != nil {
		utils.Fatal("failed to create auth service", err)
	}

	// Handlers
	authH := rest.NewAuthHandler(authSvc)

	// Route groups
	api := e.Group("/api")
	authenticatedApi := api.Group("", echojwt.WithConfig(echojwt.Config{SigningKey: []byte(cfg.JWTSecret)}))

	// Unauthenticated routes
	api.POST("/auth/login", authH.Login)

	// Authenticated routes
	authenticatedApi.POST("/", func(c echo.Context) error {
		return c.String(200, "authenticated")
	})

	log.Fatal(e.Start(cfg.ServerAddr))
}
