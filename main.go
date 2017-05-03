package main

import (
	health "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	app := cli.App("methode-article-image-set-mapper", "Maps inline image-sets from bodies of Methode articles.")
	appSystemCode := app.String(cli.StringOpt{
		Name:   "app-system-code",
		Value:  "methode-article-image-set-mapper",
		Desc:   "System Code of the application",
		EnvVar: "APP_SYSTEM_CODE",
	})
	appName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  "methode-article-image-set-mapper",
		Desc:   "Application name",
		EnvVar: "APP_NAME",
	})
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})
	logrus.Infof("methode-article-image-set-mapper is starting...\n")
	app.Action = func() {
		logrus.Infof("systemCode=%s appName=%s port=%s", *appSystemCode, *appName, *port)
		router := mux.NewRouter()
		routeProductionEndpoints(router)
		routeAdminEndpoints(router, *appSystemCode, *appName)
		if err := http.ListenAndServe(":"+*port, router); err != nil {
			logrus.Fatalf("Cound't serve http endpoints. %v\n", err)
		}
		waitForSignals()
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatalf("methode-article-image-set-mapper could not start. %v\n", err)
	}
}

func routeProductionEndpoints(router *mux.Router) {
	mapperService := newMapper()
	router.Path("/map").Handler(handlers.MethodHandler{"POST": http.HandlerFunc(mapperService.mapArticleImageSets)})
}

func routeAdminEndpoints(router *mux.Router, appSystemCode string, appName string) {
	healthService := newHealthService(&healthConfig{appSystemCode: appSystemCode, appName: appName})

	hc := health.HealthCheck{SystemCode: appSystemCode, Name: appName, Description: "Maps inline image-sets from bodies of Methode articles.", Checks: healthService.checks}
	router.Path(healthPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(health.Handler(hc))})
	router.Path(status.GTGPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.NewGoodToGoHandler(healthService.gtgCheck))})
	router.Path(status.BuildInfoPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.BuildInfoHandler)})
}

func waitForSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
