package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

// Server wraps http.Server
type Server struct {
	Server *http.Server

	gracefulShutdownTimeout int64

	logger *zap.SugaredLogger
}

// NewServer creates Server
func NewServer(server *http.Server, sugarLogger *zap.SugaredLogger, gracefulShutdownTimeout int64) *Server {
	return &Server{
		Server:                  server,
		gracefulShutdownTimeout: gracefulShutdownTimeout,
		logger:                  sugarLogger,
	}
}

// Run runs server and wait os.Interrupt signal to gracefully shutdown server
// and close all resources with closers
func (s *Server) Run(ctx context.Context, closers ...func()) {
	const fn = "Run"

	closing := false

	// run server
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && !(closing && err == http.ErrServerClosed) {
			// when shutdown here will be http.ErrServerClosed, never mind about that
			s.logger.Fatalf("%s: Server error: %v", fn, err)
		}
	}()
	s.logger.Infof("%s: Starting to listen %s...", fn, s.Server.Addr)

	// default closer which shutdowns server
	closer := func() {
		const fn = "closer"
		closing = true

		shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(s.gracefulShutdownTimeout)*time.Second)
		defer cancel()

		if err := s.Server.Shutdown(shutdownCtx); err != nil {
			s.logger.Errorf("%s: Shutdown error: %v", fn, err)
		}
	}

	s.gracefulShutdown(append(closers, closer)...)
}

// gracefulShutdown waits for interrupt signal and gracefully close all resources
func (s *Server) gracefulShutdown(closers ...func()) {
	const fn = "gracefulShutdown"

	// signal for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	s.logger.Infof("%s: Shutting down...", fn)
	for _, closer := range closers {
		closer()
	}
}
