package models

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DatabaseModel struct {
	DB       *gorm.DB
	Dialect  string
	ErrorMap func(error) interface{}
}

var DB *DatabaseModel

func Setup(db *gorm.DB) error {
	DB = &DatabaseModel{DB: db, Dialect: db.Dialector.Name()}
	switch DB.Dialect {
	case "sqlite":
		duplicate := regexp.MustCompile(`UNIQUE constraint failed: (.*)`)
		DB.ErrorMap = func(err error) interface{} {
			errText := err.Error()
			if match := duplicate.FindStringSubmatch(errText); len(match) == 2 {
				return gin.H{"error": fmt.Sprintf("Duplicate value %s", match[1])}
			}
			return gin.H{"error": errText}
		}
	case "postgres":
		duplicate := regexp.MustCompile(`ERROR: duplicate key value violates unique constraint "(.*)" \(SQLSTATE 23505\)`)
		DB.ErrorMap = func(err error) interface{} {
			errText := err.Error()
			if match := duplicate.FindStringSubmatch(errText); len(match) == 2 {
				switch match[1] {
				case "idx_books_title":
					return gin.H{"error": "Duplicate value books.title"}
				}
			}
			return gin.H{"error": errText}
		}
	default:
		DB.ErrorMap = func(err error) interface{} {
			return gin.H{"error": err.Error()}
		}
	}

	return nil
}
