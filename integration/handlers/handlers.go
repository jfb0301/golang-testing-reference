package handlers

import (
	"net/http"

	"github.com/jfb0301/golang-testing-reference/integration/db"
)

// Handler contains the handler and all its dependencies.
type Handler struct {
	bs *db.BookService
	us *db.UserService
}

// NewHandler initialises a new handler, given dependencies.
func NewHandler(bs *db.BookService, us *db.UserService) *Handler {
	return &Handler{
		bs: bs,
		us: us,
	}
}

// Index is invoked by HTTP GET /.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	resp := &Response{
		Message: "Welcome to the BookSwap service!",
		Books:   h.bs.List(),
	}
	writeResponse(w, http.StatusOK, resp)
}