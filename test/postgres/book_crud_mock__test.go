package test

import (
	"database/sql/driver"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/senomas/go-api/models"
	test_base "github.com/senomas/go-api/test/base"
	test_lib "github.com/senomas/go-api/test/lib"
	"gorm.io/driver/postgres"
)

func TestBook(t *testing.T) {
	if sqlDB, mock, err := sqlmock.New(); err != nil {
		t.Fatal("init SQLMock Error", err)
	} else {
		dialector := postgres.New(postgres.Config{
			Conn: sqlDB,
		})

		test_base.TestBookCRUD(t, dialector, mock, func(name string) {
			switch name {
			case "AutoMigrate":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`)).WithArgs("books", "BASE TABLE").WillReturnRows(sqlmock.NewRows(
					[]string{"TABLES"}))

				mock.ExpectExec(test_lib.QuoteMeta(`CREATE TABLE "books" ("id" bigserial,"title" text,"author" text,"summary" text,PRIMARY KEY ("id"))`)).WithArgs([]driver.Value{}...).WillReturnResult(driver.RowsAffected(1))

				mock.ExpectExec(test_lib.QuoteMeta(`CREATE UNIQUE INDEX IF NOT EXISTS "idx_books_title" ON "books" ("title")`)).WithArgs([]driver.Value{}...).WillReturnResult(driver.RowsAffected(1))
			case "TestBook/Finds_Empty":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books"`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(0))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 1000`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"id", "title", "author", "summary"}))
			case "TestBook/Insert_Harry_Potter_and_the_Philosopher's_Stone":
				mock.ExpectBegin()
				mock.ExpectQuery(test_lib.QuoteMeta(`INSERT INTO "books" ("title","author","summary") VALUES ($1,$2,$3) RETURNING "id"`)).WithArgs("Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			case "TestBook/Insert_Harry_Potter_and_the_Chamber_of_Secrets":
				mock.ExpectBegin()
				mock.ExpectQuery(test_lib.QuoteMeta(`INSERT INTO "books" ("title","author","summary") VALUES ($1,$2,$3) RETURNING "id"`)).WithArgs("Harry Potter and the Chamber of Secrets", "J. K. Rawling", "").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectCommit()
			case "TestBook/Finds":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books"`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(2))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 1000`)).WithArgs([]driver.Value{}...).WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", ""))
			case "TestBook/Finds_Chamber_of_Secrets":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title LIKE $1`)).WithArgs("%Chamber of Secrets%").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(1))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 LIMIT 1000`)).WithArgs("%Chamber of Secrets%").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", ""))
			case "TestBook/Finds_chamber_of_secrets":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title LIKE $1`)).WithArgs("%chamber of secrets%").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(0))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 LIMIT 1000`)).WithArgs("%chamber of secrets%").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}))
			case "TestBook/Finds_chamber_of_secrets_using_ILIKE":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title ILIKE $1`)).WithArgs("%chamber of secrets%").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(1))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title ILIKE $1 LIMIT 1000`)).WithArgs("%chamber of secrets%").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", ""))
			case "TestBook/Insert_Harry_Potter_and_Book_of_Dark_Magic":
				mock.ExpectBegin()
				mock.ExpectQuery(test_lib.QuoteMeta(`INSERT INTO "books" ("title","author","summary") VALUES ($1,$2,$3) RETURNING "id"`)).WithArgs("Harry Potter and Book of Dark Magic", "Lord Voldermort", "").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
				mock.ExpectCommit()
			case "TestBook/Finds_include_evil_book":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books"`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(3))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 1000`)).WithArgs([]driver.Value{}...).WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", "").
						AddRow(3, "Harry Potter and Book of Dark Magic", "Lord Voldermort", ""))
			case "TestBook/Finds_goods_book_only":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE NOT (author = $1)`)).WithArgs("Lord Voldermort").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(2))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE NOT (author = $1) LIMIT 1000`)).WithArgs("Lord Voldermort").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", ""))
			case "TestBook/Insert_Tintin_in_Tibet":
				mock.ExpectBegin()
				mock.ExpectQuery(test_lib.QuoteMeta(`INSERT INTO "books" ("title","author","summary") VALUES ($1,$2,$3) RETURNING "id"`)).WithArgs("Tintin in Tibet", "Herge", "").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))
				mock.ExpectCommit()
			case "TestBook/Finds_many_books":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books"`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(4))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 1000`)).WithArgs([]driver.Value{}...).WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", "").
						AddRow(3, "Harry Potter and Book of Dark Magic", "Lord Voldermort", "").
						AddRow(4, "Tintin in Tibet", "Herge", ""))
			case "TestBook/Finds_Harry_Potter_books":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title LIKE $1`)).WithArgs("%Harry Potter%").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(3))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 LIMIT 1000`)).WithArgs("%Harry Potter%").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", "").
						AddRow(3, "Harry Potter and Book of Dark Magic", "Lord Voldermort", ""))
			case "TestBook/Finds_Harry_Potter_books_from_J._K._Rawling":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title LIKE $1 AND author = $2`)).WithArgs("%Harry Potter%", "J. K. Rawling").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(2))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 AND author = $2 LIMIT 1000`)).
					WithArgs("%Harry Potter%", "J. K. Rawling").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", ""))
			case "TestBook/Delete_Evil_book":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).
					WithArgs("3").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(3, "Harry Potter and Book of Dark Magic", "Lord Voldermort", ""))
				mock.ExpectBegin()
				mock.ExpectExec(test_lib.QuoteMeta(`DELETE FROM "books" WHERE "books"."id" = $1`)).
					WithArgs(3).WillReturnResult(driver.RowsAffected(1))
				mock.ExpectCommit()
			case "TestBook/Finds_many_good_books":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books"`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(3))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 1000`)).WithArgs([]driver.Value{}...).WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", "").
						AddRow(4, "Tintin in Tibet", "Herge", ""))
			case "TestBook/Insert_Tintin_in_Jakarta":
				mock.ExpectBegin()
				mock.ExpectQuery(test_lib.QuoteMeta(`INSERT INTO "books" ("title","author","summary") VALUES ($1,$2,$3) RETURNING "id"`)).WithArgs("Tintin in Jakarta", "Herge", "").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
				mock.ExpectCommit()
			case "TestBook/Finds_tintin_books":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title LIKE $1`)).WithArgs("%Tintin%").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(2))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 LIMIT 1000`)).WithArgs("%Tintin%").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(4, "Tintin in Tibet", "Herge", "").
						AddRow(5, "Tintin in Jakarta", "Herge", ""))
			case "TestBook/Update_typo_Tintin_in_America":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).
					WithArgs("5").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(5, "Tintin in Jakarta", "Herge", ""))
				mock.ExpectBegin()
				mock.ExpectExec(test_lib.QuoteMeta(`UPDATE "books" SET "title"=$1,"author"=$2 WHERE "id" = $3`)).
					WithArgs("Tintin in America", "Herge", 5).WillReturnResult(driver.RowsAffected(1))
				mock.ExpectCommit()
			case "TestBook/Finds_updated_tintin_books":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title LIKE $1`)).WithArgs("%Tintin%").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(2))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 LIMIT 1000`)).WithArgs("%Tintin%").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(4, "Tintin in Tibet", "Herge", "").
						AddRow(5, "Tintin in America", "Herge", ""))
			case "TestBook/Finds_with_limit":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books"`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(4))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" LIMIT 2`)).WithArgs([]driver.Value{}...).WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(1, "Harry Potter and the Philosopher's Stone", "J. K. Rawling", "The boy who lived").
						AddRow(2, "Harry Potter and the Chamber of Secrets", "J. K. Rawling", ""))
			case "TestBook/Insert_Duplicate_Tintin_in_America":
				mock.ExpectBegin()
				mock.ExpectQuery(test_lib.QuoteMeta(`INSERT INTO "books" ("title","author","summary") VALUES ($1,$2,$3) RETURNING "id"`)).WithArgs("Tintin in America", "Herge", "").WillReturnError(errors.New(`ERROR: duplicate key value violates unique constraint "idx_books_title" (SQLSTATE 23505)`))
				mock.ExpectRollback()
			case "TestBook/Update_unknown_book":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).WithArgs("9999").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}))
			case "TestBook/Delete_unknown_book":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).WithArgs("9999").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}))
			case "TestBook/Get_unknown_book":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).WithArgs("9999").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}))
			case "TestBook/Update_lead_to_duplicate_record":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE id = $1 ORDER BY "books"."id" LIMIT 1`)).
					WithArgs("5").WillReturnRows(
					sqlmock.NewRows([]string{"id", "title", "author", "summary"}).
						AddRow(5, "Tintin in Jakarta", "Herge", ""))
				mock.ExpectBegin()
				mock.ExpectExec(test_lib.QuoteMeta(`UPDATE "books" SET "title"=$1,"author"=$2 WHERE "id" = $3`)).
					WithArgs("Harry Potter and the Philosopher's Stone", "Herge", 5).WillReturnError(errors.New(`ERROR: duplicate key value violates unique constraint "idx_books_title" (SQLSTATE 23505)`))
				mock.ExpectRollback()
			case "TestBook/Finds_books_id,_title_only":
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books"`)).WithArgs([]driver.Value{}...).WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(4))
				mock.ExpectQuery(test_lib.QuoteMeta(`SELECT "id","title" FROM "books" ORDER BY "id" DESC LIMIT 1000`)).WithArgs([]driver.Value{}...).WillReturnRows(
					sqlmock.NewRows([]string{"id", "title"}).
						AddRow(5, "Tintin in America").
						AddRow(4, "Tintin in Tibet").
						AddRow(2, "Harry Potter and the Chamber of Secrets").
						AddRow(1, "Harry Potter and the Philosopher's Stone"))
			default:
				log.Printf("UNKNOWN mock name '%s'", name)
			}
		})
	}
}

func TestBook_2(t *testing.T) {
	if sqlDB, mock, err := sqlmock.New(); err != nil {
		t.Fatal("init SQLMock Error", err)
	} else {
		dialector := postgres.New(postgres.Config{
			Conn: sqlDB,
		})

		ctx := test_base.NewTestContext(t, dialector, mock, func(name string) {})
		defer ctx.Close()

		t.Run("GET tintin books", func(t *testing.T) {
			defer func() {
				if err := mock.ExpectationsWereMet(); err != nil {
					ctx.Api.T.Error("ExpectationsNotMet", err)
				}
			}()
			mock.ExpectQuery(test_lib.QuoteMeta(`SELECT count(*) FROM "books" WHERE title LIKE $1`)).WithArgs("%Tintin%").WillReturnRows(sqlmock.NewRows(
				[]string{"count"}).AddRow(0))
			mock.ExpectQuery(test_lib.QuoteMeta(`SELECT * FROM "books" WHERE title LIKE $1 LIMIT 10`)).WithArgs("%Tintin%").WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "author", "summary"}))

			ctx.Api.HttpGet("/books?limit=10&query="+models.NewQuery(nil, models.NewCondition().Like("title", "Tintin"), nil).String(),
				200,
				test_base.ResponseBooks{
					Count: 0,
					Data:  []models.Book{},
				})
		})
	}
}
