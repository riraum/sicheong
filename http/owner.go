package http

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/riraum/si-cheong/db"
)

var (
	errNotAuthenticated = errors.New("not authenticated")
	errNotOwner         = errors.New("not the author of this post")
)

// ownedPost returns the post with the given ID if the logged-in author wrote it.
func (s Server) ownedPost(r *http.Request, id int) (db.Post, error) {
	author, ok, err := s.authenticated(r)
	if err != nil {
		return db.Post{}, err
	}

	if !ok {
		return db.Post{}, errNotAuthenticated
	}

	p, err := s.DB.ReadPost(id)
	if err != nil {
		return db.Post{}, fmt.Errorf("read post: %w", err)
	}

	if p.AuthorID != author.ID {
		return db.Post{}, errNotOwner
	}

	return p, nil
}

// ownerErrorStatus maps ownership errors to HTTP status codes.
func ownerErrorStatus(err error) int {
	switch {
	case errors.Is(err, errNotAuthenticated):
		return http.StatusUnauthorized
	case errors.Is(err, errNotOwner):
		return http.StatusForbidden
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
