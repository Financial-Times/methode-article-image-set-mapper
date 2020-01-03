package main

import (
	"net/http"

	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type routing struct {
	httpMappingHandler HTTPMappingHandler
	healthCheck        *HealthCheck
	router             *mux.Router
}

func newRouting(httpMappingHandler HTTPMappingHandler, healthCheck *HealthCheck) routing {
	r := routing{
		httpMappingHandler: httpMappingHandler,
		healthCheck:        healthCheck,
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
	r.router.Path(healthPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(r.healthCheck.Health())})
	r.router.Path(status.GTGPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.NewGoodToGoHandler(r.healthCheck.GTG))})
	r.router.Path(status.BuildInfoPath).Handler(handlers.MethodHandler{"GET": http.HandlerFunc(status.BuildInfoHandler)})
	r.router.Path(status.PingPath).HandlerFunc(status.PingHandler)
}

func (r routing) listenAndServe(port string) {
	err := http.ListenAndServe(":"+port, r.router)
	if err != nil {
		logrus.Fatalf("Couldn't serve http endpoints. %v\n", err)
	}
}
