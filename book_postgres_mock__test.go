package main

import (
	"database/sql/driver"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
	"github.com/stretchr/testify/suite"
)

type ResponseBooks struct {
	Data []models.Book `json:"data"`
}

type ResponseBook struct {
	Data models.Book `json:"data"`
}

type Response struct {
	Data bool `json:"data"`
}

type ResponseError struct {
	Error string `json:"error"`
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteEnv))
}

func (suite *TestSuiteEnv) SetupSuite() {
	log.Println("setup suite")
	gin.SetMode(gin.ReleaseMode)

	db, sqlDB, mock := SetupDBMock(suite)
	defer sqlDB.Close()

	mock.ExpectQuery(QuoteMeta(`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`)).WithArgs("books", "BASE TABLE").WillReturnRows(sqlmock.NewRows(
		[]string{"TABLES"}))

	mock.ExpectExec(QuoteMeta(`CREATE TABLE "books" ("id" bigserial,"title" text,"author" text,PRIMARY KEY ("id"))`)).WillReturnResult(driver.RowsAffected(1))
	db.AutoMigrate(&models.Book{})
	if err := mock.ExpectationsWereMet(); err != nil {
		suite.Fail("Mock", err)
	}

	r := gin.Default()
	SetupRoutes(r)

	suite.server = httptest.NewServer(r)
	log.Printf("Test server start %s", suite.server.URL)
}

// Running after all tests are completed
func (suite *TestSuiteEnv) TearDownSuite() {
	log.Println("teardown suite")
	suite.server.Close()
}

func (suite *TestSuiteEnv) TestBooks_FindBooks() {
	db, sqlDB, mock := SetupDBMock(suite)
	defer sqlDB.Close()
	models.DB = db

	mock.ExpectQuery(QuoteMeta(`SELECT * FROM "books" LIMIT 1000`)).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge").AddRow(200, "Sorcerer Stone", "JK Rowling"))

	suite.HttpGet("/books", ResponseBooks{
		Data: []models.Book{
			{
				ID:     100,
				Title:  "The Adventures of Tintin",
				Author: "Herge",
			},
			{
				ID:     200,
				Title:  "Sorcerer Stone",
				Author: "JK Rowling",
			},
		},
	})
	if err := mock.ExpectationsWereMet(); err != nil {
		suite.Fail("Mock", err)
	}
}

func (suite *TestSuiteEnv) TestBooks_FindBooks_WithCondition() {
	db, sqlDB, mock := SetupDBMock(suite)
	defer sqlDB.Close()
	models.DB = db

	mock.ExpectQuery(QuoteMeta(`SELECT * FROM "books" WHERE title like $1 AND author like $2 LIMIT 200 OFFSET 100`)).WithArgs("%tintin%", "JK").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge").AddRow(200, "Sorcerer Stone", "JK Rowling"))

	suite.HttpPost("/books", controllers.QueryBookInput{FindQuery: controllers.FindQuery{Offset: 100, Limit: 200}, Title_Like: "%tintin%", Author_Like: "JK"}, ResponseBooks{
		Data: []models.Book{
			{
				ID:     100,
				Title:  "The Adventures of Tintin",
				Author: "Herge",
			},
			{
				ID:     200,
				Title:  "Sorcerer Stone",
				Author: "JK Rowling",
			},
		},
	})
	if err := mock.ExpectationsWereMet(); err != nil {
		suite.Fail("Mock", err)
	}
}

func (suite *TestSuiteEnv) TestBooks_FindBook() {
	db, sqlDB, mock := SetupDBMock(suite)
	defer sqlDB.Close()
	models.DB = db

	mock.ExpectQuery(QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).WithArgs("100").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge"))

	suite.HttpGet("/books/100", ResponseBook{
		Data: models.Book{
			ID:     100,
			Title:  "The Adventures of Tintin",
			Author: "Herge",
		},
	})
	if err := mock.ExpectationsWereMet(); err != nil {
		suite.Fail("Mock", err)
	}
}

func (suite *TestSuiteEnv) TestBooks_Create() {
	db, sqlDB, mock := SetupDBMock(suite)
	defer sqlDB.Close()
	models.DB = db

	mock.ExpectBegin()
	mock.ExpectQuery(QuoteMeta(`INSERT INTO "books" ("title","author") VALUES ($1,$2) RETURNING "id"`)).WithArgs("The Adventures of Tintin", "Herge").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
	mock.ExpectCommit()

	suite.HttpPut("/books", controllers.CreateBookInput{
		Title:  "The Adventures of Tintin",
		Author: "Herge",
	}, ResponseBook{
		Data: models.Book{
			ID:     100,
			Title:  "The Adventures of Tintin",
			Author: "Herge",
		},
	})
	if err := mock.ExpectationsWereMet(); err != nil {
		suite.Fail("Mock", err)
	}
}

func (suite *TestSuiteEnv) TestBooks_Update() {
	db, sqlDB, mock := SetupDBMock(suite)
	defer sqlDB.Close()
	models.DB = db

	mock.ExpectQuery(QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).WithArgs("100").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge"))
	mock.ExpectBegin()
	mock.ExpectExec(QuoteMeta(`UPDATE "books" SET "title"=$1,"author"=$2 WHERE "id" = $3`)).WithArgs("Tintin in Tibet", "Herge", 100).WillReturnResult(driver.RowsAffected(1))
	mock.ExpectCommit()

	suite.HttpPatch("/books/100", controllers.UpdateBookInput{
		Title:  "Tintin in Tibet",
		Author: "Herge",
	}, ResponseBook{
		Data: models.Book{
			ID:     100,
			Title:  "Tintin in Tibet",
			Author: "Herge",
		},
	})
	if err := mock.ExpectationsWereMet(); err != nil {
		suite.Fail("Mock", err)
	}
}

func (suite *TestSuiteEnv) TestBooks_Delete() {
	db, sqlDB, mock := SetupDBMock(suite)
	defer sqlDB.Close()
	models.DB = db

	mock.ExpectQuery(QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).WithArgs("100").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge"))
	mock.ExpectBegin()
	mock.ExpectExec(QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).WithArgs(int64(100)).WillReturnResult(sqlmock.NewResult(100, 1))
	mock.ExpectCommit()

	suite.HttpDelete("/books/100", Response{
		Data: true,
	})
}
