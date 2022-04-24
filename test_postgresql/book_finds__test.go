package test_postgresql

import (
	"database/sql"
	"database/sql/driver"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
	"github.com/senomas/go-api/test_lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestContext struct {
	api   *test_lib.Api
	sqlDB *sql.DB
	mock  sqlmock.Sqlmock
	db    *gorm.DB
}

type Response struct {
	Data bool `json:"data"`
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
	server := httptest.NewServer(r)

	ctx := &TestContext{api: &test_lib.Api{Server: server, T: t}}

	if sqlDB, mock, err := sqlmock.New(); err != nil {
		t.Fatal("init SQLMock Error", err)
	} else {
		ctx.sqlDB = sqlDB
		ctx.mock = mock
	}

	if db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: ctx.sqlDB,
	}), &gorm.Config{}); err != nil {
		t.Fatal("Init GORM Error", err)
	} else {
		ctx.mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`)).WithArgs("books", "BASE TABLE").WillReturnRows(sqlmock.NewRows(
			[]string{"TABLES"}))

		ctx.mock.ExpectExec(test_lib.QuoteMeta(`CREATE TABLE "books" ("id" bigserial,"title" text,"author" text,"summary" text,PRIMARY KEY ("id"))`)).WithArgs([]driver.Value{}...).WillReturnResult(driver.RowsAffected(1))

		ctx.mock.ExpectExec(test_lib.QuoteMeta(`CREATE UNIQUE INDEX IF NOT EXISTS "idx_books_title" ON "books" ("title")`)).WithArgs([]driver.Value{}...).WillReturnResult(driver.RowsAffected(1))

		db.AutoMigrate(&models.Book{})
		models.DB = db
		ctx.db = db

		if err := ctx.mock.ExpectationsWereMet(); err != nil {
			ctx.api.T.Fatal("Mock", err)
		}
	}

	return ctx
}

func (ctx *TestContext) Close() {
	if err := ctx.mock.ExpectationsWereMet(); err != nil {
		ctx.api.T.Error("Mock", err)
	}
	ctx.api.Server.Close()
	ctx.sqlDB.Close()
}

func TestBook_Finds_All(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.Close()

	ctx.mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 1000`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "title", "author"}).AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling").AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling"))

	ctx.api.HttpGet("/books", 200, ResponseBooks{
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

func TestBook_Finds_LimitAndOffset(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.Close()

	ctx.mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 10 OFFSET 30`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "title", "author"}).AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling").AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling"))

	ctx.api.HttpGet("/books?offset=30&limit=10", 200, ResponseBooks{
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

func TestBook_Finds_Filter(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.Close()

	ctx.mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 AND NOT (author = $2) LIMIT 10 OFFSET 30`)).WithArgs("%Harry Potter%", "Lord Voldermort").WillReturnRows(sqlmock.NewRows(
		[]string{"id", "title", "author"}).AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling").AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling"))

	ctx.api.HttpPost("/books?offset=30&limit=10", controllers.NewQuery(nil, controllers.NewCondition().Like("title", "Harry Potter").Not(controllers.NewCondition().Equal("author", "Lord Voldermort"))), 200, ResponseBooks{
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

func TestBook_Finds_SelectTitle(t *testing.T) {
	ctx := NewTestContext(t)
	defer ctx.Close()

	ctx.mock.ExpectQuery(test_lib.QuoteMeta(`SELECT "title" FROM "books" LIMIT 10 OFFSET 30`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
		[]string{"title"}).AddRow("Harry Potter and the Philosopher's Stone").AddRow("Harry Potter and the Chamber of Secrets"))

	ctx.api.HttpPost("/books?offset=30&limit=10", controllers.NewQuery([]string{"title"}, controllers.NewCondition()), 200, ResponseBooks{
		Data: []models.Book{
			{
				ID:     0,
				Title:  "Harry Potter and the Philosopher's Stone",
				Author: "",
			},
			{
				ID:     0,
				Title:  "Harry Potter and the Chamber of Secrets",
				Author: "",
			},
		},
	})
}
