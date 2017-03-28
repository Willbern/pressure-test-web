package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dataman-Cloud/pressure-test-web/model"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/smallnest/goreq"
	"golang.org/x/net/context"
)

var c context.Context

func (api *Api) StaticHandler(ctx *gin.Context) {
	ctx.File("./static/index.html")
}

func (api *Api) Pong(ctx *gin.Context) {
	HttpOkResponse(ctx, "pong")
}

func (api *Api) ResponseJson(ctx *gin.Context) {
	_, b, errs := goreq.New().SetHeader("Connection", "Keep- Alive").Get("http://" + api.Config.DemoAppAddr + "/json").End()
	if len(errs) != 0 {
		log.Error("get json error: ", errs[0].Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, errs[0])
	}

	ctx.String(http.StatusOK, b)
	return
}

func (api *Api) Add(ctx *gin.Context) {
	var app model.App
	err := ctx.BindJSON(&app)
	if err != nil {
		log.Error("bind json error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusBadRequest, err)
	}

	_, err = api.HttpClient.POST(c, "http://"+api.Config.DemoAppAddr+"/add", nil, app, nil)
	if err != nil {
		log.Error("add app error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	HttpOkResponse(ctx, "OK")
}

func (api *Api) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("convert strin to int error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	r, err := api.HttpClient.GET(c, "http://"+api.Config.DemoAppAddr+"/get/"+idStr, nil, nil)
	if err != nil {
		log.Error("add app error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	var app model.App
	err = json.Unmarshal(r, &app)
	if err != nil {
		log.Error("unmarshal app error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	HttpOkResponse(ctx, app)
}

func (api *Api) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("convert strin to int error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	_, err = api.HttpClient.DELETE(c, "http://"+api.Config.DemoAppAddr+"/delete/"+idStr, nil, nil)
	if err != nil {
		log.Error("add app error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	HttpOkResponse(ctx, "OK")
}

func (api *Api) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")

	var app model.App
	err := ctx.BindJSON(&app)
	if err != nil {
		log.Error("update bind json error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	_, err = api.HttpClient.PUT(c, "http://"+api.Config.DemoAppAddr+"/update/"+idStr, nil, app, nil)
	if err != nil {
		log.Error("add app error: ", err.Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, err)
	}

	HttpOkResponse(ctx, "OK")
}

func HttpOkResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
	return
}

func HttpErrorResponse(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, err.Error())
	return
}
