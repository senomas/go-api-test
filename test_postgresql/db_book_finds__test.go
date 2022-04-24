package test_postgresql

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/senomas/go-api/controllers"
	"github.com/senomas/go-api/models"
	"github.com/senomas/go-api/test_lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestContextDB struct {
	api *test_lib.Api
	db  *gorm.DB
}

func NewTestContextDB(t *testing.T) *TestContextDB {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	controllers.SetupRoutes(r)
	server := httptest.NewServer(r)

	ctx := &TestContextDB{api: &test_lib.Api{Server: server, T: t}}

	if db, err := gorm.Open(postgres.Open(`host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta`), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		t.Fatal("Init GORM Error", err)
	} else {
		db.Migrator().DropTable(&models.Book{})
		db.AutoMigrate(&models.Book{})
		models.DB = db
		ctx.db = db
	}

	return ctx
}

func (ctx *TestContextDB) Close() {
	ctx.api.Server.Close()
}

func TestDBBook(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := NewTestContextDB(t)
	defer ctx.Close()

	t.Run("Finds_Empty", func(t *testing.T) {
		ctx.api.HttpGet("/books", 200, ResponseBooks{
			Data: []models.Book{},
		})
	})

	t.Run("Insert Harry Potter and the Philosopher's Stone", func(t *testing.T) {
		ctx.api.HttpPut("/books", controllers.CreateBookInput{
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
		ctx.api.HttpPut("/books", controllers.CreateBookInput{
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
		ctx.api.HttpGet("/books", 200, ResponseBooks{
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
		ctx.api.HttpPost(
			"/books",
			controllers.NewQuery(nil, controllers.NewCondition().Like("title", "Chamber of Secrets")),
			200,
			ResponseBooks{
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
		ctx.api.HttpPost(
			"/books",
			controllers.NewQuery(nil, controllers.NewCondition().Like("title", "chamber of secrets")),
			200,
			ResponseBooks{
				Data: []models.Book{},
			})
	})

	t.Run("Finds chamber of secrets using ILIKE", func(t *testing.T) {
		ctx.api.HttpPost(
			"/books",
			controllers.NewQuery(nil, controllers.NewCondition().ILike("title", "chamber of secrets")),
			200,
			ResponseBooks{
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
		ctx.api.HttpPut("/books", controllers.CreateBookInput{
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
		ctx.api.HttpGet("/books", 200, ResponseBooks{
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
		ctx.api.HttpPost("/books",
			controllers.NewQuery(nil, controllers.NewCondition().Not(controllers.NewCondition().Equal("author", "Lord Voldermort"))),
			200,
			ResponseBooks{
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
		ctx.api.HttpPut("/books", controllers.CreateBookInput{
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
		ctx.api.HttpGet("/books", 200, ResponseBooks{
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
		ctx.api.HttpPost("/books",
			controllers.NewQuery(nil, controllers.NewCondition().Like("title", "Harry Potter")),
			200,
			ResponseBooks{
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
		ctx.api.HttpPost("/books",
			controllers.NewQuery(nil, controllers.NewCondition().Like("title", "Harry Potter").Equal("author", "J. K. Rawling")),
			200,
			ResponseBooks{
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
		ctx.api.HttpDelete("/books/3",
			200,
			Response{
				Data: true,
			})
	})

	t.Run("Finds many good books", func(t *testing.T) {
		ctx.api.HttpGet("/books", 200, ResponseBooks{
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
		ctx.api.HttpPut("/books", controllers.CreateBookInput{
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
		ctx.api.HttpPost("/books",
			controllers.NewQuery(nil, controllers.NewCondition().Like("title", "Tintin")),
			200,
			ResponseBooks{
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
		ctx.api.HttpPatch("/books/5", controllers.UpdateBookInput{
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
		ctx.api.HttpPost("/books",
			controllers.NewQuery(nil, controllers.NewCondition().Like("title", "Tintin")),
			200,
			ResponseBooks{
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

	t.Run("Insert Duplicate Tintin in America", func(t *testing.T) {
		ctx.api.HttpPut("/books", controllers.CreateBookInput{
			Title:  "Tintin in America",
			Author: "Herge",
		}, 400, map[string]string{
			"error": "ERROR: duplicate key value violates unique constraint \"idx_books_title\" (SQLSTATE 23505)",
		})
	})

	t.Run("Update unknown book", func(t *testing.T) {
		ctx.api.HttpPatch("/books/9999", controllers.CreateBookInput{
			Title:  "Book of Unknown",
			Author: "John Doe",
		}, 400, map[string]string{
			"error": "record not found",
		})
	})

	t.Run("Delete unknown book", func(t *testing.T) {
		ctx.api.HttpDelete("/books/9999", 400, map[string]string{
			"error": "record not found",
		})
	})

	t.Run("Get unknown book", func(t *testing.T) {
		ctx.api.HttpGet("/books/9999", 400, map[string]string{
			"error": "record not found",
		})
	})

	t.Run("Update lead to duplicate record", func(t *testing.T) {
		ctx.api.HttpPatch("/books/5", controllers.UpdateBookInput{
			Title:  "Harry Potter and the Philosopher's Stone",
			Author: "Herge",
		}, 400, map[string]string{
			"error": "ERROR: duplicate key value violates unique constraint \"idx_books_title\" (SQLSTATE 23505)",
		})
	})

	t.Run("Finds books id, title only", func(t *testing.T) {
		ctx.api.HttpPost("/books",
			controllers.NewQuery(controllers.Fields("id", "title"), controllers.NewCondition()),
			200,
			ResponseBooks{
				Data: []models.Book{
					{
						ID:    1,
						Title: "Harry Potter and the Philosopher's Stone",
					},
					{
						ID:    2,
						Title: "Harry Potter and the Chamber of Secrets",
					},
					{
						ID:    4,
						Title: "Tintin in Tibet",
					},
					{
						ID:    5,
						Title: "Tintin in America",
					},
				},
			})
	})
}
