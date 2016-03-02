package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tantan/conf"
	"github.com/tantan/controllers"
)

var (
	date, rev string
)

func main() {
	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginIns := gin.New()
	ginIns.Use(gin.Recovery())

	ginIns.GET("/1/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"date": date, "git": rev})
	})

	//用户相关:注册,登陆

	userAPIV1 := ginIns.Group("/")
	{
		userAPIV1.POST("/users", controllers.Register)
		userAPIV1.GET("/users", controllers.ListAll)

		userAPIV1.GET("/users/:user_id/relationships", controllers.GetAllRelation)
		userAPIV1.PUT("/users/:user_id/relationships/:other_user_id", controllers.LikeUser)
	}

	ginIns.Run(fmt.Sprintf(":%d", conf.HttpPort))

}
