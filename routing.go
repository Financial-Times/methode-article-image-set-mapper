package main

import (
	health "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/httphandlers"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type routing struct {
	httpMappingHandler HTTPMappingHandler
	healthService      *healthService
	router             *mux.Router
}

func newRouting(httpMappingHandler HTTPMappingHandler, healthService *healthService) routing {
	r := routing{
		httpMappingHandler: httpMappingHandler,
		healthService:      healthService,
		router:             mux.NewRouter(),
	}
	r.routeProductionEndpoints()
	r.routeAdminEndpoints()
	return r
}

func (r routing) routeProductionEndpoints() {
	r.router.Path("/map").Handler(handlers.MethodHandler{"POST": http.HandlerFunc(r.httpMappingHandler.handle)})
}

func (r routing) routeAdminEndpoints() {
	hc := health.HealthCheck{
		SystemCode:  r.healthService.config.appSystemCode,
		Name:        r.healthService.config.appName,
		Description: "Maps inline image-sets from bodies of Methode articles.",
		Checks:      r.healthService.checks,
	}
	r.router.Path(healthPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(health.Handler(hc))})
	r.router.Path(status.GTGPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.NewGoodToGoHandler(r.healthService.gtgCheck))})
	r.router.Path(status.BuildInfoPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.BuildInfoHandler)})
	r.router.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)
}

func (r routing) listenAndServe(port string) {
	err := http.ListenAndServe(":"+port, r.router)
	if err != nil {
		logrus.Fatalf("Cound't serve http endpoints. %v\n", err)
	}
}
