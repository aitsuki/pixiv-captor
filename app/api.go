package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"aitsuki.com/pixiv-capture/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Run(port int, dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	dao, err := data.NewIllustDao(db)
	if err != nil {
		return err
	}
	repo := data.NewIllustRepository(dao)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://www.pixiv.net"},
		AllowMethods: []string{"HEAD, GET, POST"},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       12 * time.Hour,
	}))
	r.HEAD("/pixiv/:id", checkExists(repo))
	r.GET("/pixiv/:id", get(repo))
	r.POST("/pixiv", caputre(repo))
	r.GET("/pixiv", query(repo))
	return r.Run(fmt.Sprintf(":%d", port))
}

func caputre(repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		illustData := IllustData{}
		err := c.ShouldBindJSON(&illustData)
		if err != nil {
			c.Status(http.StatusBadRequest)
			log.Println(err)
			return
		}
		err = repo.Save(illustData.ToEntity())
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		c.Status(http.StatusOK)
	}
}

func checkExists(repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if repo.IsExists(id) {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusNotFound)
		}
	}
}

func get(repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		illust, err := repo.GetById(id)
		if err != nil {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, illust)
		}
	}
}

func query(repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		r18Str := c.DefaultQuery("r18", "0")
		r18, err := strconv.Atoi(r18Str)
		if err != nil {
			r18 = 0
		}

		limitStr := c.DefaultQuery("limit", "1")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 1
		}

		q, isQuery := c.GetQuery("q")

		var illusts []data.Illust
		if isQuery {
			illusts, err = repo.Search(r18, q, limit)
		} else {
			illusts, err = repo.GetRandom(r18, limit)
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, illusts)
	}
}
