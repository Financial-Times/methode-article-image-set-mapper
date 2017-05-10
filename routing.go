package main

import (
	health "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/Sirupsen/logrus"
)

type routing struct {
	httpMappingHandler    HTTPMappingHandler
	router                *mux.Router
}

func newRouting(httpMappingHandler HTTPMappingHandler, appSystemCode string, appName string) routing {
	r := routing{
		httpMappingHandler:    httpMappingHandler,
		router:                mux.NewRouter(),
	}
	r.routeProductionEndpoints()
	r.routeAdminEndpoints(appSystemCode, appName)
	return r
}

func (r routing) routeProductionEndpoints() {
	r.router.Path("/map").Handler(handlers.MethodHandler{"POST": http.HandlerFunc(r.httpMappingHandler.handle)})
}

func (r routing) routeAdminEndpoints(appSystemCode string, appName string) {
	healthService := newHealthService(&healthConfig{appSystemCode: appSystemCode, appName: appName})

	hc := health.HealthCheck{SystemCode: appSystemCode, Name: appName, Description: "Maps inline image-sets from bodies of Methode articles.", Checks: healthService.checks}
	r.router.Path(healthPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(health.Handler(hc))})
	r.router.Path(status.GTGPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.NewGoodToGoHandler(healthService.gtgCheck))})
	r.router.Path(status.BuildInfoPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.BuildInfoHandler)})
}

func (r routing) listenAndServe(port string) {
	err := http.ListenAndServe(":"+port, r.router)
	if err != nil {
		logrus.Fatalf("Cound't serve http endpoints. %v\n", err)
	}
}