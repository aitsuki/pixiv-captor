package services

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aitsuki/pixiv-captor/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Run(port int, dbPath string, logPath string, username string, password string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	logWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(logWriter)
	gin.DefaultWriter = logWriter
	repo := data.NewIllustRepository(db)
	repo.Prepare()
	if err != nil {
		return err
	}

	gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(CORS)
	r.HEAD("/pixiv/:id", checkExists(repo))
	r.GET("/pixiv/:id", getting(repo))
	r.GET("/pixiv", query(repo))
	authorize := r.Group("/", gin.BasicAuthForRealm(gin.Accounts{username: password}, "capture"))
	authorize.DELETE("/pixiv/:id", deleting(username, password, repo))
	authorize.POST("/pixiv", capture(username, password, repo))
	return r.Run(fmt.Sprintf(":%d", port))
}

var CORS = cors.New(cors.Config{
	AllowOrigins: []string{"https://www.pixiv.net"},
	AllowMethods: []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodDelete},
	AllowHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Authorization"},
	MaxAge:       12 * time.Hour,
})

func capture(username string, password string, repo *data.IllustRepository) gin.HandlerFunc {
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
		c.Status(http.StatusNoContent)
	}
}

func checkExists(repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if repo.IsExists(id) {
			c.Status(http.StatusNoContent)
		} else {
			c.Status(http.StatusNotFound)
		}
	}
}

func getting(repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		illust, err := repo.GetByID(id)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, illust)
		}
	}
}

func deleting(username string, password string, repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := repo.Delete(id)
		if err != nil {
			log.Print(err)
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusNoContent)
		}
	}
}

func query(repo *data.IllustRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		r18Str := c.DefaultQuery("r18", "0")
		r18, err := strconv.Atoi(r18Str)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
		}

		limitStr := c.DefaultQuery("limit", "1")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			log.Println(err, limit)
			c.Status(http.StatusBadRequest)
		}

		if limit > 100 {
			limit = 100
		}

		q, isQuery := c.GetQuery("q")

		var illusts []data.Illust
		if isQuery {
			illusts, err = repo.Search(r18, q, limit)
		} else {
			illusts, err = repo.GetRandom(r18, limit)
		}
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, illusts)
	}
}
