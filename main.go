package main

import (
	health "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/Sirupsen/logrus"
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
	log.SetLevel(log.InfoLevel)
	log.Infof("[Startup] methode-article-image-set-mapper is starting ")
	app.Action = func() {
		log.Infof("System code: %s, App Name: %s, Port: %s", *appSystemCode, *appName, *port)
		router := mux.NewRouter()
		routeProductionEndpoints(router)
		routeAdminEndpoints(router, *appSystemCode, *appName)
		if err := http.ListenAndServe(":"+*port, router); err != nil {
			log.Fatalf("Unable to start: %v", err)
		}
		waitForSignals()
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}

func routeProductionEndpoints(router *mux.Router) {
	mapperService := Mapper{}
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
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
