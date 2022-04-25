package test_base

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
	test_lib "github.com/senomas/go-api/test/lib"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestContext struct {
	Api       *test_lib.Api
	mock      sqlmock.Sqlmock
	initMock  func(name string)
	db        *gorm.DB
	dialector gorm.Dialector
}

type Response struct {
	Data bool `json:"data"`
}

type ResponseBooks struct {
	Count int64         `json:"count"`
	Data  []models.Book `json:"data"`
}

type ResponseBook struct {
	Data models.Book `json:"data"`
}

func NewTestContext(t *testing.T, dialector gorm.Dialector, mock sqlmock.Sqlmock, initMock func(name string)) *TestContext {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	controllers.SetupRoutes(r)
	server := httptest.NewServer(r)

	ctx := &TestContext{Api: &test_lib.Api{Server: server, T: t}, dialector: dialector, mock: mock, initMock: initMock}

	var config *gorm.Config
	if os.Getenv("DEBUG") != "" {
		config = &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					LogLevel:                  logger.Info,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				}),
		}
	} else {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}
	if db, err := gorm.Open(ctx.dialector, config); err != nil {
		t.Fatal("Init GORM Error", err)
	} else {
		if mock == nil {
			db.Migrator().DropTable(&models.Book{})
		}
		models.Setup(db)
		ctx.db = db
	}

	defer ctx.startMock("AutoMigrate")()
	ctx.db.AutoMigrate(&models.Book{})

	return ctx
}

func (ctx *TestContext) Close() {
	ctx.Api.Server.Close()
}

func (ctx *TestContext) startMock(name string) func() {
	if ctx.mock != nil {
		ctx.initMock(name)
		return func() {
			if err := ctx.mock.ExpectationsWereMet(); err != nil {
				ctx.Api.T.Error("ExpectationsNotMet", err)
			}
		}
	} else {
		return func() {}
	}
}

func TestBookCRUD(t *testing.T, dialector gorm.Dialector, mock sqlmock.Sqlmock, initMock func(name string)) {
	if testing.Short() {
		t.Skip()
	}
	ctx := NewTestContext(t, dialector, mock, initMock)
	defer ctx.Close()

	t.Run("Finds_Empty", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpGet("/books", 200, ResponseBooks{
			Count: 0,
			Data:  []models.Book{},
		})
	})

	t.Run("Insert Harry Potter and the Philosopher's Stone", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPut("/books", controllers.CreateBookInput{
			Title:   "Harry Potter and the Philosopher's Stone",
			Author:  "J. K. Rawling",
			Summary: "The boy who lived",
		}, 200, ResponseBook{
			Data: models.Book{
				ID:      1,
				Title:   "Harry Potter and the Philosopher's Stone",
				Author:  "J. K. Rawling",
				Summary: "The boy who lived",
			},
		})
	})

	t.Run("Insert Harry Potter and the Chamber of Secrets", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPut("/books", controllers.CreateBookInput{
			Title:  "Harry Potter and the Chamber of Secrets",
			Author: "J. K. Rawling",
		}, 200, ResponseBook{
			Data: models.Book{
				ID:     2,
				Title:  "Harry Potter and the Chamber of Secrets",
				Author: "J. K. Rawling",
			},
		})
	})

	t.Run("Finds", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpGet("/books", 200, ResponseBooks{
			Count: 2,
			Data: []models.Book{
				{
					ID:      1,
					Title:   "Harry Potter and the Philosopher's Stone",
					Author:  "J. K. Rawling",
					Summary: "The boy who lived",
				},
				{
					ID:     2,
					Title:  "Harry Potter and the Chamber of Secrets",
					Author: "J. K. Rawling",
				},
			},
		})
	})

	t.Run("Finds Chamber of Secrets", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost(
			"/books",
			models.NewQuery(nil, models.NewCondition().Like("title", "Chamber of Secrets"), nil),
			200,
			ResponseBooks{
				Count: 1,
				Data: []models.Book{
					{
						ID:     2,
						Title:  "Harry Potter and the Chamber of Secrets",
						Author: "J. K. Rawling",
					},
				},
			})
	})

	t.Run("Finds chamber of secrets", func(t *testing.T) {
		if ctx.dialector.Name() == "sqlite" || ctx.dialector.Name() == "mysql" {
			t.Skipf("SKIP: same behaviour as ilike in %v", ctx.dialector.Name())
		}

		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost(
			"/books",
			models.NewQuery(nil, models.NewCondition().Like("title", "chamber of secrets"), nil),
			200,
			ResponseBooks{
				Count: 0,
				Data:  []models.Book{},
			})
	})

	t.Run("Finds chamber of secrets using ILIKE", func(t *testing.T) {
		if ctx.dialector.Name() == "sqlite" || ctx.dialector.Name() == "mysql" {
			t.Skipf("SKIP: same behaviour as ilike in %v", ctx.dialector.Name())
		}

		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost(
			"/books",
			models.NewQuery(nil, models.NewCondition().ILike("title", "chamber of secrets"), nil),
			200,
			ResponseBooks{
				Count: 1,
				Data: []models.Book{
					{
						ID:     2,
						Title:  "Harry Potter and the Chamber of Secrets",
						Author: "J. K. Rawling",
					},
				},
			})
	})

	t.Run("Insert Harry Potter and Book of Dark Magic", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPut("/books", controllers.CreateBookInput{
			Title:  "Harry Potter and Book of Dark Magic",
			Author: "Lord Voldermort",
		}, 200, ResponseBook{
			Data: models.Book{
				ID:     3,
				Title:  "Harry Potter and Book of Dark Magic",
				Author: "Lord Voldermort",
			},
		})
	})

	t.Run("Finds include evil book", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpGet("/books", 200, ResponseBooks{
			Count: 3,
			Data: []models.Book{
				{
					ID:      1,
					Title:   "Harry Potter and the Philosopher's Stone",
					Author:  "J. K. Rawling",
					Summary: "The boy who lived",
				},
				{
					ID:     2,
					Title:  "Harry Potter and the Chamber of Secrets",
					Author: "J. K. Rawling",
				},
				{
					ID:     3,
					Title:  "Harry Potter and Book of Dark Magic",
					Author: "Lord Voldermort",
				},
			},
		})
	})

	t.Run("Finds goods book only", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost("/books",
			models.NewQuery(nil, models.NewCondition().Not(models.NewCondition().Equal("author", "Lord Voldermort")), nil),
			200,
			ResponseBooks{
				Count: 2,
				Data: []models.Book{
					{
						ID:      1,
						Title:   "Harry Potter and the Philosopher's Stone",
						Author:  "J. K. Rawling",
						Summary: "The boy who lived",
					},
					{
						ID:     2,
						Title:  "Harry Potter and the Chamber of Secrets",
						Author: "J. K. Rawling",
					},
				},
			})
	})

	t.Run("Insert Tintin in Tibet", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPut("/books", controllers.CreateBookInput{
			Title:  "Tintin in Tibet",
			Author: "Herge",
		}, 200, ResponseBook{
			Data: models.Book{
				ID:     4,
				Title:  "Tintin in Tibet",
				Author: "Herge",
			},
		})
	})

	t.Run("Finds many books", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpGet("/books", 200, ResponseBooks{
			Count: 4,
			Data: []models.Book{
				{
					ID:      1,
					Title:   "Harry Potter and the Philosopher's Stone",
					Author:  "J. K. Rawling",
					Summary: "The boy who lived",
				},
				{
					ID:     2,
					Title:  "Harry Potter and the Chamber of Secrets",
					Author: "J. K. Rawling",
				},
				{
					ID:     3,
					Title:  "Harry Potter and Book of Dark Magic",
					Author: "Lord Voldermort",
				},
				{
					ID:     4,
					Title:  "Tintin in Tibet",
					Author: "Herge",
				},
			},
		})
	})

	t.Run("Finds Harry Potter books", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost("/books",
			models.NewQuery(nil, models.NewCondition().Like("title", "Harry Potter"), nil),
			200,
			ResponseBooks{
				Count: 3,
				Data: []models.Book{
					{
						ID:      1,
						Title:   "Harry Potter and the Philosopher's Stone",
						Author:  "J. K. Rawling",
						Summary: "The boy who lived",
					},
					{
						ID:     2,
						Title:  "Harry Potter and the Chamber of Secrets",
						Author: "J. K. Rawling",
					},
					{
						ID:     3,
						Title:  "Harry Potter and Book of Dark Magic",
						Author: "Lord Voldermort",
					},
				},
			})
	})

	t.Run("Finds Harry Potter books from J. K. Rawling", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost("/books",
			models.NewQuery(nil, models.NewCondition().Like("title", "Harry Potter").Equal("author", "J. K. Rawling"), nil),
			200,
			ResponseBooks{
				Count: 2,
				Data: []models.Book{
					{
						ID:      1,
						Title:   "Harry Potter and the Philosopher's Stone",
						Author:  "J. K. Rawling",
						Summary: "The boy who lived",
					},
					{
						ID:     2,
						Title:  "Harry Potter and the Chamber of Secrets",
						Author: "J. K. Rawling",
					},
				},
			})
	})

	t.Run("Delete Evil book", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpDelete("/books/3",
			200,
			Response{
				Data: true,
			})
	})

	t.Run("Finds many good books", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpGet("/books", 200, ResponseBooks{
			Count: 3,
			Data: []models.Book{
				{
					ID:      1,
					Title:   "Harry Potter and the Philosopher's Stone",
					Author:  "J. K. Rawling",
					Summary: "The boy who lived",
				},
				{
					ID:     2,
					Title:  "Harry Potter and the Chamber of Secrets",
					Author: "J. K. Rawling",
				},
				{
					ID:     4,
					Title:  "Tintin in Tibet",
					Author: "Herge",
				},
			},
		})
	})

	t.Run("Insert Tintin in Jakarta", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPut("/books", controllers.CreateBookInput{
			Title:  "Tintin in Jakarta",
			Author: "Herge",
		}, 200, ResponseBook{
			Data: models.Book{
				ID:     5,
				Title:  "Tintin in Jakarta",
				Author: "Herge",
			},
		})
	})

	t.Run("Finds tintin books", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost("/books",
			models.NewQuery(nil, models.NewCondition().Like("title", "Tintin"), nil),
			200,
			ResponseBooks{
				Count: 2,
				Data: []models.Book{
					{
						ID:     4,
						Title:  "Tintin in Tibet",
						Author: "Herge",
					},
					{
						ID:     5,
						Title:  "Tintin in Jakarta",
						Author: "Herge",
					},
				},
			})
	})

	t.Run("Update typo Tintin in America", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPatch("/books/5", controllers.UpdateBookInput{
			Title:  "Tintin in America",
			Author: "Herge",
		}, 200, ResponseBook{
			Data: models.Book{
				ID:     5,
				Title:  "Tintin in America",
				Author: "Herge",
			},
		})
	})

	t.Run("Finds updated tintin books", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost("/books",
			models.NewQuery(nil, models.NewCondition().Like("title", "Tintin"), nil),
			200,
			ResponseBooks{
				Count: 2,
				Data: []models.Book{
					{
						ID:     4,
						Title:  "Tintin in Tibet",
						Author: "Herge",
					},
					{
						ID:     5,
						Title:  "Tintin in America",
						Author: "Herge",
					},
				},
			})
	})

	t.Run("Finds with limit", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpGet("/books?limit=2", 200, ResponseBooks{
			Count: 4,
			Data: []models.Book{
				{
					ID:      1,
					Title:   "Harry Potter and the Philosopher's Stone",
					Author:  "J. K. Rawling",
					Summary: "The boy who lived",
				},
				{
					ID:     2,
					Title:  "Harry Potter and the Chamber of Secrets",
					Author: "J. K. Rawling",
				},
			},
		})
	})

	t.Run("Insert Duplicate Tintin in America", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPut("/books", controllers.CreateBookInput{
			Title:  "Tintin in America",
			Author: "Herge",
		}, 400, map[string]string{
			"error": "Duplicate value books.title",
		})
	})

	t.Run("Update unknown book", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPatch("/books/9999", controllers.CreateBookInput{
			Title:  "Book of Unknown",
			Author: "John Doe",
		}, 400, map[string]string{
			"error": "record not found",
		})
	})

	t.Run("Delete unknown book", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpDelete("/books/9999", 400, map[string]string{
			"error": "record not found",
		})
	})

	t.Run("Get unknown book", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpGet("/books/9999", 400, map[string]string{
			"error": "record not found",
		})
	})

	t.Run("Update lead to duplicate record", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPatch("/books/5", controllers.UpdateBookInput{
			Title:  "Harry Potter and the Philosopher's Stone",
			Author: "Herge",
		}, 400, map[string]string{
			"error": "Duplicate value books.title",
		})
	})

	t.Run("Finds books id, title only", func(t *testing.T) {
		defer ctx.startMock(t.Name())()

		ctx.Api.HttpPost("/books",
			models.NewQuery(models.Fields("id", "title"), nil, &models.QueryOrderBy{Field: "id", Desc: true}),
			200,
			ResponseBooks{
				Count: 4,
				Data: []models.Book{
					{
						ID:    5,
						Title: "Tintin in America",
					},
					{
						ID:    4,
						Title: "Tintin in Tibet",
					},
					{
						ID:    2,
						Title: "Harry Potter and the Chamber of Secrets",
					},
					{
						ID:    1,
						Title: "Harry Potter and the Philosopher's Stone",
					},
				},
			})
	})
}
