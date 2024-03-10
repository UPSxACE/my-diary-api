package server

import (
	"net/http"

	"github.com/UPSxACE/my-diary-api/db"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostNoteBody struct {
	Title string `json:"title" validate:"required,max=254"`
	Content string `json:"content" validate:"required,max=131070"`
	ContentRaw string `json:"contentRaw" validate:"required,max=131070"`
}


func (s *Server) postNotesRoute(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)

	jwtId := claims.UserId

	// Read body
	noteBody := &PostNoteBody{}

	if err := c.Bind(noteBody); err != nil {
		return echo.ErrBadRequest
	}

	// Validate fields
	err := s.validator.Struct(noteBody)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		if len(errs) > 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"field": errs[0].Field()})
		}
	}

	// Save
	params := db.CreateNoteParams{AuthorID: int32(jwtId), Title: noteBody.Title, Content: noteBody.Content, ContentRaw: noteBody.ContentRaw}

	id, err := s.Queries.CreateNote(c.Request().Context(), params)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, id)
}