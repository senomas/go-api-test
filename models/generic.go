package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (db *DatabaseModel) Finds(c *gin.Context, model interface{}, data interface{}) {
	tx := db.DB.Model(model)
	var query Query
	if c.Request.Method == "POST" {
		if err := c.ShouldBindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if str := c.Query("query"); str != "" {
		if err := json.Unmarshal([]byte(str), &query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	where, params := query.Condition.Apply("", []any{})
	tx.Where(where, params...)
	if query.OrderBy.Field != "" {
		tx.Order(clause.OrderByColumn{Column: clause.Column{Name: query.OrderBy.Field}, Desc: query.OrderBy.Desc})
	}

	var count int64
	if err := tx.Count(&count).Error; err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorMap(err))
		return
	}

	if query.Select != nil {
		tx = tx.Select(query.Select)
	}

	if str := c.Query("offset"); str != "" {
		if i, err := strconv.Atoi(str); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Offset error": err.Error()})
		} else {
			tx = tx.Offset(i)
		}
	}
	if str := c.Query("limit"); str != "" {
		if i, err := strconv.Atoi(str); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Limit error": err.Error()})
		} else {
			tx = tx.Limit(i)
		}
	} else {
		tx = tx.Limit(1000)
	}

	if err := tx.Find(data).Error; err != nil {
		c.JSON(http.StatusBadRequest, DB.ErrorMap(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "data": data})
}

func (db *DatabaseModel) Find(c *gin.Context, data interface{}) {
	if err := db.DB.Where("id = ?", c.Param("id")).First(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorMap(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (db *DatabaseModel) Create(c *gin.Context, data interface{}) {
	if err := db.DB.Create(data).Error; err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorMap(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (db *DatabaseModel) Delete(c *gin.Context, model interface{}) {
	if tx := db.DB.Where("id = ?", c.Param("id")).First(model); tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tx.Error.Error()})
		return
	} else if tx.RowsAffected != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid RowsAffected %v", tx.RowsAffected)})
		return
	}

	if err := db.DB.Delete(model).Error; err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorMap(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (db *DatabaseModel) Update(c *gin.Context, data interface{}, applyInput func()) {
	if err := db.DB.Where("id = ?", c.Param("id")).First(data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	applyInput()

	if err := db.DB.Updates(data).Error; err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorMap(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
