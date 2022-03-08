package rest

import (
	"fmt"
	"ravxcheckout/src/adapter/config"
	v1 "ravxcheckout/src/adapter/rest/router/v1"

	"github.com/gin-gonic/gin"
)

func InitRest() {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	cfg := config.GetConfig()

	routerV1 := engine.Group("/v1")
	{
		v1.InitializeRouter(routerV1)
	}

	engine.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	restAddr := fmt.Sprintf("0.0.0.0:%d", cfg.RestPort)

	go func() {
		err := engine.Run(restAddr)

		if err != nil {
			panic(err.Error())
		}
	}()

}
