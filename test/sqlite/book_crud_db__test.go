package test

import (
	"testing"

	test_base "github.com/senomas/go-api/test/base"
	"gorm.io/driver/sqlite"
)

func TestBookDB(t *testing.T) {
	dialector := sqlite.Open("file::memory:?cache=shared")

	test_base.TestBookCRUD(t, dialector, nil, func(name string) {})
}
