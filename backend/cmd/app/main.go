package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"book-tracker/config"
	mariadbRepo "book-tracker/internal/repository/mariadb"
	"book-tracker/internal/rest"
	"book-tracker/services/auth"
	"book-tracker/utils"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("failed to get config", err)
	}

	utils.SetupLogger(cfg.LogLevel)

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
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			msg := fmt.Sprintf("request timed out on route %s", c.Path())
			utils.Fatal(msg, err)
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

	userRepo := mariadbRepo.NewUserRepository(dbConn)

	authSvc := auth.NewAuthService(userRepo)

	rest.NewAuthHandler(e, authSvc)

	log.Fatal(e.Start(cfg.ServerAddr))
}
