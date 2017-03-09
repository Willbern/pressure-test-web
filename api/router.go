package api

import (
	"github.com/Dataman-Cloud/pressure-test-web/config"
	"github.com/Dataman-Cloud/pressure-test-web/httpclient"

	//"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Config     *config.Config
	HttpClient *httpclient.Client
}

func (api *Api) ApiRouter() *gin.Engine {
	router := gin.New()

	router.GET("/ping", api.Pong)
	router.POST("/add", api.Add)
	router.GET("/get/:id", api.Get)
	router.DELETE("/delete/:id", api.Delete)
	router.PUT("/update/:id", api.Update)

	//router.GET("/static", api.StaticHandler)
	router.Static("/ui", "./static")
	router.StaticFile("/pic", "./static/snowbording.jpg")

	return router
}
