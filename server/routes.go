package server

import "github.com/labstack/echo/v4"

func (s *Server) setRoutes(devMode bool) {
	// Public Routes
	s.router.GET("/ping", func(c echo.Context) error {
		return c.JSON(200,"pong");
	})

	// Guest Routes
	routeIndexGuest := s.router.Group("/", s.guestMiddleware)
	routeIndexGuest.POST("login", s.postLoginRoute)
	routeIndexGuest.POST("register", s.postRegisterRoute)

	// Private Routes
	routeIndexPrivate := s.router.Group("/", s.jwtMiddleware)
	routeIndexPrivate.GET("profile", s.getProfileRoute)
	if(devMode){
		routeIndexPrivate.POST("blacklist-token", s.postBlacklistTokenRoute)
		routeIndexPrivate.GET("test-token", s.getTestTokenRoute)
	}
	
	

	// Moderation Routes
}
