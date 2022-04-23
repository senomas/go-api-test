package main

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
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

type TestSuiteEnv struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteEnv))
}

func (suite *TestSuiteEnv) SetupSuite() {
	log.Println("setup suite")
	gin.SetMode(gin.ReleaseMode)
}

// Running after all tests are completed
func (suite *TestSuiteEnv) TearDownSuite() {
	log.Println("teardown suite")
}

func setupDBMock(suite *TestSuiteEnv) (*gorm.DB, sqlmock.Sqlmock) {
	var err error
	var sqlDB *sql.DB
	var mock sqlmock.Sqlmock
	sqlDB, mock, err = sqlmock.New()
	if err != nil {
		suite.Errorf(err, "Failed to open mock sql db")
	}
	var db *gorm.DB
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		suite.Errorf(err, "Failed to setup mockup gorm")
	}
	return db, mock
}

func (suite *TestSuiteEnv) marshal(v any) string {
	var str string
	if bb, err := json.MarshalIndent(v, "", "\t"); err != nil {
		suite.Error(err)
	} else {
		str = string(bb)
	}
	return str
}

func (suite *TestSuiteEnv) Test_AutoMigrate_Books() {
	db, mock := setupDBMock(suite)

	mock.ExpectQuery("SHOW TABLES FROM `` WHERE `Tables_in_` = \\?").WithArgs("books").WillReturnRows(sqlmock.NewRows(
		[]string{"TABLES"}))

	mock.ExpectExec("CREATE TABLE `books` \\(`id` int unsigned AUTO_INCREMENT,`title` varchar\\(255\\),`author` varchar\\(255\\) , PRIMARY KEY \\(`id`\\)\\)").WillReturnResult(driver.RowsAffected(1))
	db.AutoMigrate(&models.Book{})
}

func (suite *TestSuiteEnv) TestBooks_FindBooks() {
	a := suite.Assert()

	db, mock := setupDBMock(suite)
	models.DB = db

	mock.ExpectQuery(`SELECT \* FROM "books"`).WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge").AddRow(200, "Sorcerer Stone", "JK Rowling"))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	controllers.FindBooks(ctx)
	a.Equal(200, recorder.Code, "Response Code")
	var res ResponseBooks
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		a.Error(err)
	}
	a.Equal(suite.marshal(ResponseBooks{
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
	}), suite.marshal(res), recorder.Body.String())
}

func (suite *TestSuiteEnv) TestBooks_FindBook() {
	a := suite.Assert()

	db, mock := setupDBMock(suite)
	models.DB = db

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE id = \$1 ORDER BY "books"."id" LIMIT 1`).WithArgs("100").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge"))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "100"})

	controllers.FindBook(ctx)
	a.Equal(200, recorder.Code, "Response Code")
	var res ResponseBook
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		a.Error(err)
	}
	a.Equal(suite.marshal(ResponseBook{
		Data: models.Book{
			ID:     100,
			Title:  "The Adventures of Tintin",
			Author: "Herge",
		},
	}), suite.marshal(res), recorder.Body.String())
}

func (suite *TestSuiteEnv) TestBooks_FindBook_NotFound() {
	a := suite.Assert()

	db, mock := setupDBMock(suite)
	models.DB = db

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE id = \$1 ORDER BY "books"."id" LIMIT 1`).WithArgs("100").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "100"})

	controllers.FindBook(ctx)
	a.Equal(400, recorder.Code, "Response Code")
	var res ResponseError
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		a.Error(err)
	}
	a.Equal(suite.marshal(ResponseError{
		Error: "Record not found!",
	}), suite.marshal(res), recorder.Body.String())
}
func (suite *TestSuiteEnv) TestBooks_Create() {
	a := suite.Assert()

	db, mock := setupDBMock(suite)
	models.DB = db

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "books" \("title","author"\) VALUES \(\$1,\$2\) RETURNING "id"`).WithArgs("The Adventures of Tintin", "Herge").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))
	mock.ExpectCommit()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	bookBytes, _ := json.Marshal(controllers.CreateBookInput{
		Title:  "The Adventures of Tintin",
		Author: "Herge",
	})
	ctx.Request = httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(bookBytes))

	controllers.CreateBook(ctx)
	a.Equal(200, recorder.Code, "Response Code")
	var res ResponseBook
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		a.Error(err)
	}
	a.Equal(suite.marshal(ResponseBook{
		Data: models.Book{
			ID:     100,
			Title:  "The Adventures of Tintin",
			Author: "Herge",
		},
	}), suite.marshal(res), recorder.Body.String())
}

func (suite *TestSuiteEnv) TestBooks_Update() {
	a := suite.Assert()

	db, mock := setupDBMock(suite)
	models.DB = db

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE id = \$1 ORDER BY "books"."id" LIMIT 1`).WithArgs("100").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge"))
	// mock.ExpectExec(`UPDATExxx`).WithArgs("Herge", "Tintin in Tibet", 100).WillReturnResult(driver.RowsAffected(1))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	bookBytes, _ := json.Marshal(controllers.UpdateBookInput{
		Title:  "Tintin in Tibet",
		Author: "Herge",
	})
	ctx.Request = httptest.NewRequest(http.MethodPatch, "/books/100", bytes.NewBuffer(bookBytes))
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "100"})

	controllers.UpdateBook(ctx)
	a.Equal(200, recorder.Code, "Response Code")
	var res ResponseBook
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		a.Error(err)
	}
	a.Equal(suite.marshal(ResponseBook{
		Data: models.Book{
			ID:     100,
			Title:  "Tintin in Tibet",
			Author: "Herge",
		},
	}), suite.marshal(res), recorder.Body.String())
}

func (suite *TestSuiteEnv) TestBooks_Delete() {
	a := suite.Assert()

	db, mock := setupDBMock(suite)
	models.DB = db

	mock.ExpectQuery(`SELECT \* FROM "books" WHERE id = \$1 ORDER BY "books"."id" LIMIT 1`).WithArgs("100").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(100, "The Adventures of Tintin", "Herge"))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "books" WHERE "books"."id" = \$1`).WithArgs(int64(100)).WillReturnResult(sqlmock.NewResult(100, 1))
	mock.ExpectCommit()

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/books/100", nil)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: "100"})

	controllers.DeleteBook(ctx)
	a.Equal(200, recorder.Code, "Response Code")
	var res Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		a.Error(err)
	}
	a.Equal(suite.marshal(Response{
		Data: true,
	}), suite.marshal(res), recorder.Body.String())
}
