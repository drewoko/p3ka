package core

import (
	"log"
	"strconv"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func Web(db *DataBase, config *Config) {

	defer func() {
		if r := recover(); r != nil {
			log.Print("Recovering WEB Instance", r)
		}
	}()

	ginInst := gin.Default()

	if config.Dev {
		ginInst.Use(static.Serve("/", static.LocalFile(config.Static, true)))
	} else {
		ginInst.Use(static.Serve("/", BinaryFileSystem("static/dist")))
		ginInst.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	ginInst.NoRoute(redirect)

	api := ginInst.Group("/api")
	{
		api.GET("/random", func(c *gin.Context) {

			filter := c.Query("filter")
			filterInt, err := strconv.Atoi(filter)

			if err != nil {
				c.JSON(400, gin.H{})
				return
			}

			c.JSON(200, db.getRandom(config.HttpResponseLimit, filterInt))
		})

		api.GET("/top", func(c *gin.Context) {
			source := c.Query("source")

			if source == "" {
				c.JSON(200, db.getTop(500, config.ExcludedUsers))
			} else {
				c.JSON(200, db.getTopUsersBySource(500, source, config.ExcludedUsers))
			}
		})

		api.GET("/last", func(c *gin.Context) {

			start := c.Query("start")
			filter := c.Query("filter")

			startInt, err := strconv.Atoi(start)
			filterInt, err2 := strconv.Atoi(filter)

			if err != nil || err2 != nil {
				c.JSON(400, gin.H{})
				return
			}

			c.JSON(200, db.GetLast(config.HttpResponseLimit, startInt, filterInt))

		})

		userGroup := api.Group("/user")

		userGroup.GET("/", func(c *gin.Context) {

			start := c.Query("start")
			user := c.Query("user")

			i, err := strconv.Atoi(start)

			if err != nil {
				c.JSON(400, gin.H{})
			} else {
				c.JSON(200, db.GetLastUser(config.HttpResponseLimit, i, user))
			}
		})

		userGroup.GET("/id", func(c *gin.Context) {

			start := c.Query("start")
			id := c.Query("id")

			startInt, err := strconv.Atoi(start)
			if err != nil {
				c.JSON(400, gin.H{})
				return
			}

			idInt, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(400, gin.H{})
				return
			}

			c.JSON(200, db.GetLastUserById(config.HttpResponseLimit, startInt, idInt))
		})
	}

	ginInst.Run(":" + config.Port)
}

func redirect(c *gin.Context) {
	c.Redirect(301, "/index.html")
}
