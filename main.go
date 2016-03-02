package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tantan/conf"
	"github.com/tantan/controllers"
	"github.com/tantan/logger"
)

var (
	date, rev string
)

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		dur := time.Since(start) / time.Millisecond

		requestLogData := controllers.GetRequestLogDataFromContext(c)
		logInfo := map[string]interface{}{
			"code":   c.Writer.Status(),
			"dur":    dur,
			"remote": c.ClientIP(),
			"url":    c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"method": c.Request.Method,
			"data":   requestLogData,
		}

		if c.Writer.Status() >= 400 {
			logger.RequestLogger.Error(logInfo)
		} else {
			logger.RequestLogger.Info(logInfo)
		}

	}
}

func main() {
	if conf.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginIns := gin.New()
	ginIns.Use(gin.Recovery())
	ginIns.Use(requestLogger())

	ginIns.GET("/1/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"date": date, "git": rev})
	})

	//用户相关:注册,登陆...

	userAPIV1 := ginIns.Group("/")
	{
		userAPIV1.POST("/users", controllers.Register)
		userAPIV1.GET("/users", controllers.ListAll)

		userAPIV1.GET("/users/:user_id/relationships", controllers.GetAllRelation)
		userAPIV1.PUT("/users/:user_id/relationships/:other_user_id", controllers.LikeUser)
	}

	err := ginIns.Run(fmt.Sprintf(":%d", conf.HttpPort))
	if err != nil {
		fmt.Println("gin start err:" + err.Error())
	}

}
