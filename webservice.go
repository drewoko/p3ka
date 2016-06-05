package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"strconv"
)

func web(db *DataBase, config *Config) {

	ginInst := gin.Default()

	ginInst.Use(static.Serve("/", static.LocalFile(config.Static, true)))

	ginInst.GET("/api/random", func(c *gin.Context) {
		c.JSON(200, db.get_random(39))
	})

	ginInst.GET("/api/top", func(c *gin.Context) {
		c.JSON(200, db.get_top_users(200))
	})

	ginInst.GET("/api/last", func(c *gin.Context) {

		start := c.Query("start")

		i, err := strconv.Atoi(start)

		if(err != nil) {
			c.JSON(400, gin.H{})
		} else {
			c.JSON(200, db.get_last(37, i))
		}
	})


	ginInst.GET("/api/user", func(c *gin.Context) {

		start := c.Query("start")
		user := c.Query("user")

		i, err := strconv.Atoi(start)

		if(err != nil) {
			c.JSON(400, gin.H{})
		} else {
			c.JSON(200, db.get_last_user(37, i, user))
		}
	})


	ginInst.Run(":"+config.Port)
}
