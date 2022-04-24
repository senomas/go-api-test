package test_sqlite

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
	"github.com/senomas/go-api/test_lib"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestContext struct {
	api *test_lib.Api
}

type ResponseBooks struct {
	Data []models.Book `json:"data"`
}

type ResponseBook struct {
	Data models.Book `json:"data"`
}

func NewTestContext(t *testing.T) *TestContext {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	controllers.SetupRoutes(r)

	if db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{}); err != nil {
		t.Fatal("SQLite Error", err)
	} else {
		db.AutoMigrate(&models.Book{})
		models.DB = db
	}

	server := httptest.NewServer(r)
	return &TestContext{api: &test_lib.Api{Server: server, T: t}}
}

func (ctx *TestContext) Close() {
	ctx.api.Server.Close()
}

func TestBooks(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.Close()

	ctx.api.HttpGet("/books", ResponseBooks{
		Data: []models.Book{},
	})

	ctx.api.HttpPut("/books", controllers.CreateBookInput{
		Title:  "The Adventures of Tintin",
		Author: "Herge",
	}, ResponseBook{
		Data: models.Book{
			ID:     1,
			Title:  "The Adventures of Tintin",
			Author: "Herge",
		},
	})

	ctx.api.HttpGet("/books", ResponseBooks{
		Data: []models.Book{
			{
				ID:     1,
				Title:  "The Adventures of Tintin",
				Author: "Herge",
			},
		},
	})
}
