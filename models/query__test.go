package models

import (
	"encoding/json"
	"log"
	"testing"

	test_lib "github.com/senomas/go-api/test/lib"
	"github.com/stretchr/testify/assert"
)

func TestBook_Simple(t *testing.T) {
	query := NewQuery([]string{"id", "title"}, NewCondition().Like("title", "harry potter"), nil)

	var bytes []byte
	if bb, err := json.MarshalIndent(query, "", "  "); err != nil {
		log.Fatal(err)
	} else {
		bytes = bb
	}

	value := NewQuery(nil, nil, nil)
	if err := json.Unmarshal(bytes, &value); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, test_lib.Marshal(t, query), test_lib.Marshal(t, value))

	where, params := value.Condition.Apply("", []any{})
	assert.Equal(t, "title LIKE ?", where, string(bytes))
	assert.Equal(t, test_lib.Marshal(t, []any{"%harry potter%"}), test_lib.Marshal(t, params))
}

func TestBook_Query(t *testing.T) {
	query := NewQuery([]string{"title"}, NewCondition().Like("title", "harry potter"), nil)

	var bytes []byte
	if bb, err := json.MarshalIndent(query, "", "  "); err != nil {
		log.Fatal(err)
	} else {
		bytes = bb
	}

	value := NewQuery(nil, nil, nil)
	if err := json.Unmarshal(bytes, &value); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, test_lib.Marshal(t, query), test_lib.Marshal(t, value))
}

func TestBook_Not(t *testing.T) {
	query := NewQuery([]string{"id", "title"}, NewCondition().Like("title", "harry potter").Not(NewCondition().Equal("author", "Lord Voldermort")), nil)

	var bytes []byte
	if bb, err := json.MarshalIndent(query, "", "  "); err != nil {
		log.Fatal(err)
	} else {
		bytes = bb
	}

	value := NewQuery(nil, nil, nil)
	if err := json.Unmarshal(bytes, &value); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, test_lib.Marshal(t, query), test_lib.Marshal(t, value))

	where, params := value.Condition.Apply("", []any{})
	assert.Equal(t, "title LIKE ? AND NOT (author = ?)", where, string(bytes))
	assert.Equal(t, test_lib.Marshal(t, []any{"%harry potter%", "Lord Voldermort"}), test_lib.Marshal(t, params))
}

func TestBook_Or(t *testing.T) {
	query := NewQuery([]string{"id", "title"}, NewCondition().Like("title", "harry potter").Or(NewCondition().Equal("author", "Lord Voldermort").Equal("author", "Tom Malvolo Riddle")), nil)

	var bytes []byte
	if bb, err := json.MarshalIndent(query, "", "  "); err != nil {
		log.Fatal(err)
	} else {
		bytes = bb
	}

	value := NewQuery(nil, nil, nil)
	if err := json.Unmarshal(bytes, &value); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, test_lib.Marshal(t, query), test_lib.Marshal(t, value))

	where, params := value.Condition.Apply("", []any{})
	assert.Equal(t, "title LIKE ? AND (author = ? OR author = ?)", where, string(bytes))
	assert.Equal(t, test_lib.Marshal(t, []any{"%harry potter%", "Lord Voldermort", "Tom Malvolo Riddle"}), test_lib.Marshal(t, params))
}

func TestBook_Fail_LikeNotString(t *testing.T) {
	bytes := []byte(`{
		"o": "AND",
		"e": [
			{
				"o": "LIKE",
				"f": "title",
				"v": 1000
			},
			{
				"o": "OR",
				"e": [
					{
						"o": "=",
						"f": "author",
						"v": "Lord Voldermort"
					},
					{
						"o": "=",
						"f": "author",
						"v": "Tom Malvolo Riddle"
					}
				]
			}
		]
	}`)

	value := NewCondition()
	assert.ErrorContains(t, json.Unmarshal(bytes, &value), "UNSUPPORTED TYPE VALUE 1000")
}

func TestBook_Fail_EqObject(t *testing.T) {
	bytes := []byte(`{
		"o": "AND",
		"e": [
			{
				"o": "=",
				"f": "title",
				"v": {
					"low": 1,
					"high": 100
				}
			},
			{
				"o": "OR",
				"e": [
					{
						"o": "=",
						"f": "author",
						"v": "Lord Voldermort"
					},
					{
						"o": "=",
						"f": "author",
						"v": "Tom Malvolo Riddle"
					}
				]
			}
		]
	}`)

	value := NewCondition()
	assert.ErrorContains(t, json.Unmarshal(bytes, &value), "UNSUPPORTED TYPE VALUE")
}

func TestBook_Fail_Not(t *testing.T) {
	bytes := []byte(`{
		"o": "AND",
		"e": [
			{
				"o": "NOT",
				"e": [
					{
						"o": "=",
						"f": "author",
						"v": "Lord Voldermort"
					},
					{
						"o": "=",
						"f": "author",
						"v": "Tom Malvolo Riddle"
					}
				]
			}
		]
	}`)

	value := NewCondition()
	assert.ErrorContains(t, json.Unmarshal(bytes, &value), "UNSUPPORTED EXPRESSION NOT")
}
