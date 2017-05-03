package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	health "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"net/http"
)

type routing struct {
	router *mux.Router
}

func newRouting(mapperService ImageSetMapper, appSystemCode string, appName string) routing {
	r := routing{
		router: mux.NewRouter(),
	}
	r.routeProductionEndpoints(mapperService)
	r.routeAdminEndpoints(appSystemCode, appName)
	return r
}

func (r routing) routeProductionEndpoints(mapperService ImageSetMapper) {
	httpMappingHandler := newHttpMappingHandler(mapperService)
	r.router.Path("/map").Handler(handlers.MethodHandler{"POST": http.HandlerFunc(httpMappingHandler.handle)})
}

func (r routing) routeAdminEndpoints(appSystemCode string, appName string) {
	healthService := newHealthService(&healthConfig{appSystemCode: appSystemCode, appName: appName})

	hc := health.HealthCheck{SystemCode: appSystemCode, Name: appName, Description: "Maps inline image-sets from bodies of Methode articles.", Checks: healthService.checks}
	r.router.Path(healthPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(health.Handler(hc))})
	r.router.Path(status.GTGPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.NewGoodToGoHandler(healthService.gtgCheck))})
	r.router.Path(status.BuildInfoPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.BuildInfoHandler)})
}

func (r routing) listenAndServe(port string) error {
	err := http.ListenAndServe(":"+port, r.router)
	if err != nil {
		return err
	}
	return nil
}
