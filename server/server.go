package server

import (
	"context"

	"github.com/UPSxACE/my-diary-api/db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	router *echo.Echo
	db *pgx.Conn
	dbContext context.Context
	Queries *db.Queries 
	tokenBlacklist sessionRevokeList
	jwtConfig echojwt.Config
	validator *validator.Validate // use a single instance of Validate, it caches struct info
}

func NewServer(devMode bool) *Server {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1323", "http://localhost:3000"},
		AllowHeaders: []string{"*"},
		// echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept
	}))

	
	server := &Server{router:e}

	server.setupValidator();
	server.setupDatabase(devMode)
	server.upgradeDatabase(devMode)
	server.setupJwt()
	server.setRoutes()

	return server
}

func (s *Server) Start(address string) error{
	defer s.db.Close(s.dbContext)
	return s.router.Start(address)
}