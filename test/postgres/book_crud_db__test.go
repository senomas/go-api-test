package test

import (
	"testing"

	test_base "github.com/senomas/go-api/test/base"
	"gorm.io/driver/postgres"
)

func TestBookDB(t *testing.T) {
	dialector := postgres.Open(`host=localhost user=demo password=password dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Jakarta`)

	test_base.TestBookCRUD(t, dialector, nil, func(name string) {})
}
