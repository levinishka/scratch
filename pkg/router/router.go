package router

import (
	"net/http/pprof"

	"github.com/causelovem/scratch/pkg/metrics"
	"github.com/gorilla/mux"
)

// NewRouter creates new mux router
func NewRouter(strictSlash bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(strictSlash)
	// always use prometheus metrics middleware
	router.Use(metrics.PrometheusMiddleware)
	return router
}

// NewRouterWithPprof creates new mux router and register pprof handlers
// path for all pprof handlers has /debug/pprof/ prefix
func NewRouterWithPprof(strictSlash bool) *mux.Router {
	router := NewRouter(strictSlash)
	addPprof(router)

	return router
}

// addPprof adds pprof handlers to router
func addPprof(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}
