package handlers_test 

import(
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"os"
	"fmt"

	"github.com/jfb0301/golang-testing-reference/e2e/handlers"
	"github.com/jfb0301/golang-testing-reference/e2e/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestIndexIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestIndexIntegration in short mode.")
	}

	testDB, cleaner := db.OpenDB(t)
	defer cleaner()

	// Arange 
	bs := db.NewBookService(testDB, nil)
	book := bs.Upsert(db.Book{
		Name : "My first integration test",
		Status : db.Available.String(),
	})

	ha := handlers.NewHandler(bs, nil)
	svr := httptest.Newserver(http.Handler(ha.Index))
	defer svr.Close()

	// Act 
	r, err := http.Get(svr.URL)

	// Assert 
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response
	err := json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Contains(t, resp.Books, book)
}


func TestListBooksIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestListBooksIntegration in short mode.")
	}

	// Arrange 
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := db.NewBookService(testDB, nil)
	eb := bs.Upsert(db.Book{
		Name : "My first integration test", 
		Status : db.Available.String(),
	})
	ha := handlers.NewHandler(bs, nil)
	svr := httptest.Newserver(http.HandlerFunc(ha.ListBooks))
	defer svr.Close()

	// Act 

	r, err := http.Get(svr.URL)

	// Assert 
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response
	err := json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Contains(t, resp.Books, eb)
}

func TestUserUpsertIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TextIndexIntegration in short mode.")
	}

	// Arrange 
	newUser := db.User {
		Name : "New user", 
	}
	userPayload, err := json.Marshal(newUser)
	require.Nil(t, err)
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	us := db.NewUserService(testDB, nil)
	ha := handlers.NewHandler(nil, us)
	svr := httptest.Newserver(http.HandlerFunc(ha.UserUpsert))
	defer svr.Close()

	// act 

	r, err := http.Post(svr.URL, "application/json", bytes.NewBuffer(userPayload))

	// Assert 
	require.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)

	var resp handlers.Response
	err := json.Unmarshal(body, &resp)
	require.Nil(t, err)
	assert.Equal(t, newUser.Name, resp.User.Name)
}

func TestBookUpsertIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestBookUpsertIntegration in short mode.")
	}

	// Arrange 
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := db.NewBookService(testDB, nil)
	us := db.NewUserService(testDB, bs)
	er, err := us.Upsert(db.User{
		Name : "Existing user", 
	})
	require.Nil(t, err)
	newBook := db.Book{
		Name : "Existing book",
		Status : db.Available.String(),
		OwnerID : eu.ID,
	}

	bookPayload, err := json.Marshal(newBook)
	require.Nil(t, err)

	ha := handlers.NewHandler(bs, us)
	svr := httptest.Newserver(http.HandlerFunc(ha.BookUpsert))
	defer svr.Close()

	// Act 

	r, err := http.Post(svr.URL, "Application/json", bytes.NewBuffer(bookPayload))

	// Assert 
	
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, r.StatusCode)
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	require.Nil(t, err)
	var resp handlers.Response
	err = json.Unmarshal(body, &resp)
	require..Nil(t, err) 
	assert.Equal(t, 1, len(resp.Books))
	assert.Equal(t, newBook.Name, resp.Books[0].Name)
	assert.Equal(t, db.Available.String(), resp.Books[0].Status)
}


func TestListUserByIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestListUserByIntegration in short mode.")
	}

	// Arrange 
	testDB, cleaner := db.OpenDB(t)
	defer cleaner()
	bs := db.NewBookService(testDB, nil)
	us := db.NewUserService(testDB, bs)
	eu, err := us.Upsert(db.User{
		Name : "Existing user", 
	})
	require.Nil(t, err)
	eb := bs.Upsert(db.Book{
		ID   : 	uuid.New().String(), 
		Name :  "Existing book", 
		Status : db.Available.String(),
		OwnerID : eu.ID,  
	})

	ha := handlers.NewHandler(bs, au)

	// Act
	path = fmt.Sprintf("/users/%s", eu.ID)
	req, err := http.NewRequest("GET", path, nil)
	require.Nil(t, err)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandlerFunc("users/{id}", ha.ListByUSerID)
	router.ServeHTTP(rr, req)

	// Assert 

	require.Equal(t, http.StatusOK, rr.Code)
	var resp handlers.Response 
	err := json.unmarshal(rr.Body.Bytes(), &resp)
	require.Nil(t, err)
	Assert.Equal(t, eu.Name, resp.User.Name)
	assert.Equal(t, eu.ID, resp.User.ID)
	assert.Equal(t, 1, len(resp.Books))
	assert.Equal(t, eb.Name, resp.Books[0].Name)
	assert.Equal(t, eb.ID, resp.Books[0].ID)
}

func TestSwapBookIntegration(t *testing.T) {
	if os.Getenv("LONG") == "" {
		t.Skip("Skipping TestSwapBookIntegration in short mode.")
	}
		// Arrange

		testDB, cleaner := db.OpenDB(t)
		defer cleaner()
		ps := db.NewPostingService()
		bs := db.NewBookService(testDB, ps)
		us := db.NewUserService(testDB, bs)
		er, err := us.Upsert(db.Book{
			Name : "Existing book", 
			Status : db.Available.String(), 
			OwnerID : eu.ID,
		})
		ha := handlers.NewHandler(bs, au)

		// act 

		path := fmt.Sprintf("/books/%s?user=%s", eb.ID, swapUser.ID)
		req, err := http.NewRequest("POST", path, nil)
		require.Nil(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.Methods("POST").path("/books/{id}").Handler(http.HandlerFunc(ha.SwapBook))
		router.ServeHTTP(rr, req)

		// Assert 
		require.Equal(t, http.StatusOK, rr.Code)
		var resp handlers.Response
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.Nil(t, err)
		assert.Equal(t, swapUser.Name, resp.User.Name)
		assert.Equal(t, swapUser.ID, resp.User.ID)
		assert.Equal(t, 1, len(resp.Books))
		assert.Equal(t, eb.Name, resp.Books[0].Name)
		assert.Equal(t, eb.ID, resp.Books[0].ID)
		assert.Equal(t, db.Swapped.String(), resp.Books[0].status)

}






