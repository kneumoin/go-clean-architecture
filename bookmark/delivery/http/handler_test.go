package http

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/go-clean-architecture/auth"
	"github.com/zhashkevych/go-clean-architecture/bookmark"
	"github.com/zhashkevych/go-clean-architecture/bookmark/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	testUser := &auth.User{
		Username: "testuser",
		Password: "testpass",
	}

	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.BookmarkUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	inp := &createInput{
		URL:   "testurl",
		Title: "testtitle",
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc.On("CreateBookmark", testUser, inp.URL, inp.Title).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bookmarks", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGet(t *testing.T) {
	testUser := &auth.User{
		Username: "testuser",
		Password: "testpass",
	}

	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.BookmarkUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	bms := make([]*bookmark.Bookmark, 5)
	for i := 0; i < 5; i++ {
		bms[i] = &bookmark.Bookmark{
			ID:    uuid.New(),
			URL:   "url",
			Title: "title",
		}
	}

	uc.On("GetBookmarks", testUser).Return(bms, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/bookmarks", nil)
	r.ServeHTTP(w, req)

	expectedOut := &getResponse{Bookmarks: toBookmarks(bms)}

	expectedOutBody, err := json.Marshal(expectedOut)
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedOutBody), w.Body.String())
}

func TestDelete(t *testing.T) {
	testUser := &auth.User{
		Username: "testuser",
		Password: "testpass",
	}

	r := gin.Default()
	group := r.Group("/api", func(c *gin.Context) {
		c.Set(auth.CtxUserKey, testUser)
	})

	uc := new(usecase.BookmarkUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	inp := &deleteInput{
		ID: uuid.New(),
	}

	body, err := json.Marshal(inp)
	assert.NoError(t, err)

	uc.On("DeleteBookmark", testUser, inp.ID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/bookmarks", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
