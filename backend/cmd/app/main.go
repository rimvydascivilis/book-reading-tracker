package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rimvydascivilis/book-tracker/backend/config"
	mariadbRepo "github.com/rimvydascivilis/book-tracker/backend/internal/repository/mariadb"
	"github.com/rimvydascivilis/book-tracker/backend/internal/rest"
	"github.com/rimvydascivilis/book-tracker/backend/services/auth"
	"github.com/rimvydascivilis/book-tracker/backend/services/user"
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
				"uri":    values.URI,
				"status": values.Status,
			})

			return nil
		},
	}))

	// Repositories
	userRepo := mariadbRepo.NewUserRepository(dbConn)

	// Services
	googleOauth2Svc, err := auth.NewGoogleOAuth2Service()
	if err != nil {
		utils.Fatal("failed to create Google OAuth2 service", err)
	}
	jwtSvc := auth.NewJWTService(cfg.JWTSecret, userRepo)
	userSvc := user.NewUserService(userRepo)
	authSvc := auth.NewAuthService(userSvc, googleOauth2Svc, jwtSvc)

	// Handlers
	authH := rest.NewAuthHandler(authSvc)

	// Route groups
	api := e.Group("/api")
	authenticatedApi := api.Group("", echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
		ContextKey: "user",
	}))

	// Unauthenticated routes
	api.POST("/auth/login", authH.Login)

	// Authenticated routes
	authenticatedApi.POST("/", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.String(500, "failed to get user from context")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.String(500, "failed to get claims from token")
		}
		userID, ok := claims["sub"].(float64)
		if !ok {
			return c.String(500, "failed to get user ID from claims")
		}

		utils.Info("authenticated user", map[string]interface{}{
			"userID": userID,
		})

		return c.String(200, "authenticated")
	})

	log.Fatal(e.Start(cfg.ServerAddr))
}
