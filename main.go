package main

import (
	"flag"
	"net/http"

	"github.com/Dataman-Cloud/pressure-test-web/api"
	"github.com/Dataman-Cloud/pressure-test-web/config"
	"github.com/Dataman-Cloud/pressure-test-web/httpclient"

	log "github.com/Sirupsen/logrus"
)

var (
	envFile = flag.String("config", "./deploy/env_file", "")
)

func main() {
	flag.Parse()

	conf := config.InitConfig(*envFile)

	client, _ := httpclient.NewClient(nil, nil)
	api := &api.Api{
		Config:     conf,
		HttpClient: client,
	}

	server := &http.Server{
		Addr:           conf.DemoAddr,
		Handler:        api.ApiRouter(),
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("can't start server: ", err)
	}
}
