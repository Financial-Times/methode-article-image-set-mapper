package main

import (
	health "github.com/Financial-Times/go-fthealth/v1_1"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type routing struct {
	messageToNativeMapper MessageToNativeMapper
	mapperService ImageSetMapper
	router *mux.Router
}

func newRouting(messageToNativeMapper MessageToNativeMapper, mapperService ImageSetMapper, appSystemCode string, appName string) routing {
	r := routing{
		messageToNativeMapper: messageToNativeMapper,
		mapperService: mapperService,
		router: mux.NewRouter(),
	}
	r.routeProductionEndpoints()
	r.routeAdminEndpoints(appSystemCode, appName)
	return r
}

func (r routing) routeProductionEndpoints() {
	httpMappingHandler := newHTTPMappingHandler(r.messageToNativeMapper, r.mapperService)
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
	return http.ListenAndServe(":"+port, r.router)
}
