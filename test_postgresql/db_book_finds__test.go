package test_postgresql

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
	"github.com/senomas/go-api/test_lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestContextDB struct {
	api   *test_lib.Api
	sqlDB *sql.DB
	db    *gorm.DB
}

func NewTestContextDB(t *testing.T) *TestContextDB {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	controllers.SetupRoutes(r)
	server := httptest.NewServer(r)

	ctx := &TestContextDB{api: &test_lib.Api{Server: server, T: t}}

	if db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: ctx.sqlDB,
	}), &gorm.Config{}); err != nil {
		t.Fatal("Init GORM Error", err)
	} else {
		db.AutoMigrate(&models.Book{})
		models.DB = db
		ctx.db = db
	}

	return ctx
}

func (ctx *TestContextDB) Close() {
	ctx.api.Server.Close()
	ctx.sqlDB.Close()
}

func TestDBBook(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := NewTestContextDB(t)
	defer ctx.Close()

	ctx.api.HttpGet("/books", ResponseBooks{
		Data: []models.Book{
			{
				ID:     1,
				Title:  "Harry Potter and the Philosopher's Stone",
				Author: "J. K. Rawling",
			},
			{
				ID:     2,
				Title:  "Harry Potter and the Chamber of Secrets",
				Author: "J. K. Rawling",
			},
		},
	})
}
