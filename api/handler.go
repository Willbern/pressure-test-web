package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
	_, b, errs := goreq.New().Get("http://" + api.Config.DemoAppAddr + "/json").End()
	if len(errs) != 0 {
		log.Error("get json error: ", errs[0].Error())
		HttpErrorResponse(ctx, http.StatusServiceUnavailable, errs[0])
	}

	ctx.String(http.StatusOK, b)
	return
}

func setHeader(request *http.Request, owner string, name string) {
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Connection", "Keep-Alive")
	request.Close = true
}

func (api *Api) ResponseJsonLocal(ctx *gin.Context) {
	request, err := http.NewRequest("GET", "http://"+api.Config.DemoAppAddr+"/json", nil)
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Connection", "Keep-Alive")
	request.Close = true
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = errors.New("response not 200")
		log.Warn(err)
		return
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	ctx.String(http.StatusOK, string(result))
}

func (api *Api) ResponseJsonNotKeepAlive(ctx *gin.Context) {
	request, err := http.NewRequest("GET", "http://"+api.Config.DemoAppAddr+"/json", nil)
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = errors.New("response not 200")
		log.Warn(err)
		return
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	ctx.String(http.StatusOK, string(result))
}

func (api *Api) ResponseJsonKeepAlive(ctx *gin.Context) {
	request, err := http.NewRequest("GET", "http://"+api.Config.DemoAppAddr+"/json", nil)
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Connection", "Keep-Alive")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = errors.New("response not 200")
		log.Warn(err)
		return
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	ctx.String(http.StatusOK, string(result))
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

	request, err := http.NewRequest("GET", "http://"+api.Config.DemoAppAddr+"/get/"+idStr, nil)
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = errors.New("response not 200")
		log.Warn(err)
		return
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var app model.App
	err = json.Unmarshal(result, &app)
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
