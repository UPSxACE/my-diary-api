package server

import "github.com/labstack/echo/v4"
type Pagination struct {
	TotalRecords int    `json:"total_records"`
	PageSize     int    `json:"page_size"`
	Cursor       string `json:"cursor"`
}

func (s *Server) setRoutes(devMode bool) {
	// Public Routes
	// - Index
	s.router.GET("/ping", func(c echo.Context) error {
		return c.JSON(200,"pong");
	})

	// Guest Routes
	routeIndexGuest := s.router.Group("/", s.guestMiddleware)
	// - Auth
	routeIndexGuest.POST("login", s.postLoginRoute)
	routeIndexGuest.POST("register", s.postRegisterRoute)

	// Private Routes
	routeIndexPrivate := s.router.Group("/", s.jwtMiddleware)
	// - Index
	routeIndexPrivate.GET("profile", s.getProfileRoute)
	// - Test
	if(devMode){
		routeIndexPrivate.POST("blacklist-token", s.postBlacklistTokenRoute)
		routeIndexPrivate.GET("test-token", s.getTestTokenRoute)
	}
	// - Notes
	routeNotePrivate := s.router.Group("/notes", s.jwtMiddleware)
	routeNotePrivate.POST("", s.postNotesRoute)

	// Moderation Routes
}
