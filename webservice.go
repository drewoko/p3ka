package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"strconv"
)

func web(db *DataBase, config *Config) {

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovering", r)
		}
	}()

	ginInst := gin.Default()

	ginInst.Use(static.Serve("/", static.LocalFile(config.Static, true)))
	ginInst.NoRoute(redirect)

	api := ginInst.Group("/api")
	{
		api.GET("/random", func(c *gin.Context) {
			c.JSON(200, db.get_random(39))
		})

		api.GET("/top", func(c *gin.Context) {
			c.JSON(200, db.get_top_users(200, config.ExcludedUsers))
		})

		api.GET("/last", func(c *gin.Context) {

			start := c.Query("start")

			i, err := strconv.Atoi(start)

			if(err != nil) {
				c.JSON(400, gin.H{})
			} else {
				c.JSON(200, db.get_last(37, i))
			}
		})

		api.GET("/user", func(c *gin.Context) {

			start := c.Query("start")
			user := c.Query("user")

			i, err := strconv.Atoi(start)

			if(err != nil) {
				c.JSON(400, gin.H{})
			} else {
				c.JSON(200, db.get_last_user(37, i, user))
			}
		})
	}

	ginInst.Run(":"+config.Port)
}

func redirect(c *gin.Context) {
	c.Redirect(301, "/index.html")
}