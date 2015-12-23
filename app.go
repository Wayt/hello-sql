package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gotoolz/env"
	"net/http"
	"strconv"
)

type Entry struct {
	ID    int64  `json:"id"`
	Value string `json:"value" form:"user" binding:"required`
}

func main() {

	db, err := sql.Open(env.GetDefault("SQL_DRIVER", "mysql"), env.GetDefault("SQL_DSN", "root@tcp(localhost:3306)/hello_sql"))
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/entries/:id", func(c *gin.Context) {

		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		e := new(Entry)
		if err := db.QueryRow("SELECT id, value FROM entries WHERE id = ?", id).
			Scan(&e.ID, &e.Value); err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, e)
	})

	r.POST("/entries", func(c *gin.Context) {

		e := new(Entry)
		if err := c.Bind(e); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		res, err := db.Exec("INSERT INTO entries (`value`) VALUES (?)", e.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		e.ID = id
		c.JSON(http.StatusCreated, e)
	})

	r.PUT("/entries/:id", func(c *gin.Context) {

		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		e := new(Entry)
		if err := c.Bind(e); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		e.ID = id

		_, err := db.Exec("UPDATE entries SET `value` = ? WHERE `id` = ?", e.Value, e.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, e)
	})

	r.DELETE("/entries/:id", func(c *gin.Context) {

		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		_, err := db.Exec("DELETE FROM entries WHERE `id` = ?", id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		c.AbortWithStatus(http.StatusOK)
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
