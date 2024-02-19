package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/UPSxACE/my-diary-api/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const TOKEN_DURATION = time.Hour * 24

type jwtCustomClaims struct {
	UserId      string `json:"userId"`
	Permissions string `json:"permissions"`
	jwt.RegisteredClaims
}

type userId = string
type revokeTime = time.Time
type tokenRevokeList = map[userId]revokeTime

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	devFlag := flag.Bool("dev", false, "Run server on developer mode")
	flag.Parse()

	USERNAME := os.Getenv("POSTGRES_USERNAME")
	PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	HOST := os.Getenv("POSTGRES_HOST")
	DATABASE := os.Getenv("POSTGRES_DATABASE")
	DATABASE_DEV := os.Getenv("POSTGRES_DATABASE_DEV")

	var connectionString string
	if *devFlag {
		connectionString = fmt.Sprintf("postgres://%v:%v@%v/%v", USERNAME, PASSWORD, HOST, DATABASE_DEV)
	} else {
		connectionString = fmt.Sprintf("postgres://%v:%v@%v/%v", USERNAME, PASSWORD, HOST, DATABASE)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	println(queries)

	// list all notes
	notes, err := queries.ListNotes(ctx)
	if err != nil {
		return err
	}
	log.Println(notes)
	// list all notes
	users, err := queries.ListUser(ctx)
	if err != nil {
		return err
	}
	log.Println(users)

	// Create e instance
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1323", "http://localhost:3000"},
		AllowHeaders: []string{"*"},
		// echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept
	}))

	SECRET := os.Getenv("JWT_SECRET")

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(SECRET),
	}

	var tokenBlacklist tokenRevokeList = map[string]time.Time{}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)
	jwtMiddlewareWrapper := func(next echo.HandlerFunc) echo.HandlerFunc {
		blacklistWrapper := func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*jwtCustomClaims)

			id := claims.UserId
			issuedAt := claims.IssuedAt

			revokeTime := tokenBlacklist[id]
			if !revokeTime.IsZero() {
				revokedAtMs := revokeTime.UnixMilli()
				issuedAtMs := issuedAt.UnixMilli()
				nowMs := time.Now().UnixMilli()
				tokenDurationMs := TOKEN_DURATION.Milliseconds()

				if issuedAtMs < revokedAtMs {
					return echo.ErrUnauthorized
				}
				if nowMs > revokedAtMs+tokenDurationMs {
					tokenBlacklist[id] = time.Time{};
				}
			}
			
			return next(c);
		}

		return jwtMiddleware(blacklistWrapper)
	}
	guestMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			hasToken := c.Request().Header.Get("Authorization") != ""

			if hasToken {
				return echo.ErrForbidden
			}

			return next(c)
		}
	}

	routeIndexGuest := e.Group("/", guestMiddleware)
	routeIndexGuest.POST("get-token", func(c echo.Context) error {
		type Login struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		login := &Login{}

		if err := c.Bind(login); err != nil {
			return echo.ErrBadRequest
		}

		if login.Username != "abc" || login.Password != "abc" {
			return echo.ErrUnauthorized
		}

		// Set claims
		claims := jwtCustomClaims{
			"1",
			"3",
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_DURATION)),
				IssuedAt: jwt.NewNumericDate(time.Now()),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		signedJwt, err := token.SignedString([]byte(SECRET))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{"token": signedJwt})
	})

	routeIndexPrivate := e.Group("/", jwtMiddlewareWrapper)
	routeIndexPrivate.POST("blacklist-token", func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*jwtCustomClaims)

		tokenBlacklist[claims.UserId] = time.Now()
		return c.NoContent(http.StatusOK)
	})
	routeIndexPrivate.GET("test-token", func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*jwtCustomClaims)

		id := claims.UserId
		permissions := claims.Permissions
		expiresAt := claims.ExpiresAt

		return c.JSON(http.StatusOK, echo.Map{
			"id":          id,
			"permissions": permissions,
			"expiresAt":   expiresAt,
		})
	})

	return e.Start(":1323")
}
