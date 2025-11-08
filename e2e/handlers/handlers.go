package handlers

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"

	"github.com/jfb0301/golang-testing-reference/e2e/db"
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
	books, err := h.bs.List()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error : err.Error(),
		}) 
		return 
	}

	// Send HTTP status & hardcoded message 
	resp := &Response{
		Message : "Welcome to the BookSwap service!", 
		Books : books, 
	}
	writeResponse(w, http.StatusOK, resp)
}

func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.b.List()
	if err != nil {
		writeResponse(w http.StatusInternalServerError, &Response{
			Error : err.Error(),
		})
		return 
	}

	writeResponse(w, http.StatusOK, &Response{
		Books: books, 
	})
}

func(h *Handler) UserUpsert(w http.ResponseWriter, r *http.Request) {
	// read the request body 
	body, err := readRequestBody(r)
	// Hadnle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error : fmt.Errorf("Invalid user body:%v", err).Error(),
		})
		return 
	}

	// Call the repository method corresponding to the operation 
	user, err := h.us.Upsert(user)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error : err.Error(), 
		})
		return 
	}
	
	// Send an HTTP success status & the return value from the repo 
	writeResponse(w, http.StatusOK, &Response{
		User : &user, 
	})
}

// ListByUSerID is invoked by HTTP GET /users/{id}

func (h *Handler) ListByUSerID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	user, books, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error : err.Error(), 
		})
		return 
	}

	// Send an HTTP success status & the return value from the repo 
	writeResponse(w, http.StatusOK, &Response{
		User: user, 
		Books: books,
	})
}


func (h *Handler) SwapBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["id"]
	userID := r.URL.Query().Get("user")
	if err := h.us.Exists(userID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error : err.Error(), 
		})
		return 
	}

	_, err := h.bs.SwapBook(bookID, userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error : err.Error(),
		})
		return 
	}

	user, books, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: err.Error(), 
		})
		return 
	}

	writeResponse(w, http.StatusOK, &Response{
		User: user, 
		Books: books, 
	})
}


// Book Upsert is invoked by the HTTP POST /books 

func(h *Handler) BookUpsert(w http.ResponseWriter, r *http.Request) {
	// read the request body 
	body, err := readRequestBody(r)

	// Handle any error and write an error HTTP status & response 
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: fmt.Errorf("Invalid book body:%v", err).Error(), 
		})
		return 
	}

	// Initialize a book to unmarshal request body into 
	var book db.Book 
	if err := json.unmarshal(body, &book); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response{
			Error : fmt.Errorf("invalid book body:%v", err).Error(), 
		})
		return 
	}

	if err := h.us.Exists(book.OwnerID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(), 
		})
		return 
	} 

	// Call the repository method corresponding to the operation 
	udpatedBook := h.bs.Upsert(book)
	// Send and HTTP success status & the return value from the repo 
	writeResponse(w, http.StatusOK, &Response{
		Books : []db.Book{udpatedBook}, 
	})
}

// readRequestBody: Helper method that allows to read a request body and return any errors 

func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1045876))
	if err != nil {
		return []byte{}, err
	}
	if err := r.body.Close(); err != nil {
		return []byte{}, err
	}
	return body, err
}