package core

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
)

func Web(db *DataBase, config *Config) {

	defer func() {
		if r := recover(); r != nil {
			log.Print("Recovering", r)
		}
	}()

	ginInst := gin.Default()

	if config.Dev {
		ginInst.Use(static.Serve("/", static.LocalFile(config.Static, true)))
	} else {
		ginInst.Use(static.Serve("/", BinaryFileSystem("static/dist")))
	}

	ginInst.NoRoute(redirect)

	api := ginInst.Group("/api")
	{
		api.GET("/random", func(c *gin.Context) {
			c.JSON(200, db.getRandom(config.HttpResponseLimit))
		})

		api.GET("/top", func(c *gin.Context) {
			c.JSON(200, db.getTopUsers(500, config.ExcludedUsers))
		})

		api.GET("/last", func(c *gin.Context) {

			start := c.Query("start")

			i, err := strconv.Atoi(start)

			if(err != nil) {
				c.JSON(400, gin.H{})
			} else {
				c.JSON(200, db.GetLast(config.HttpResponseLimit, i))
			}
		})

		userGroup := api.Group("/user");

		userGroup.GET("/", func(c *gin.Context) {

			start := c.Query("start")
			user := c.Query("user")

			i, err := strconv.Atoi(start)

			if(err != nil) {
				c.JSON(400, gin.H{})
			} else {
				c.JSON(200, db.GetLastUser(config.HttpResponseLimit, i, user))
			}
		});

		userGroup.GET("/id", func(c *gin.Context) {

			start := c.Query("start")
			id := c.Query("id")

			startInt, err := strconv.Atoi(start)
			if(err != nil) {
				c.JSON(400, gin.H{})
				return
			}

			idInt, err := strconv.Atoi(id)
			if(err != nil) {
				c.JSON(400, gin.H{})
				return
			}

			c.JSON(200, db.GetLastUserById(config.HttpResponseLimit, startInt, idInt))
		})
	}

	ginInst.Run(":"+config.Port)
}

func redirect(c *gin.Context) {
	c.Redirect(301, "/index.html")
}