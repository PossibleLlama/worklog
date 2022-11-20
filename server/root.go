package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const (
	PORT = 8080

	PATH = "/worklog"
)

var (
	httpRouter *mux.Router = mux.NewRouter().StrictSlash(true)
)

func Execute() {
	startServer()
}

func startServer() {
	fmt.Printf("server starting on port %d\n", PORT)

	httpRouter.NotFoundHandler = http.HandlerFunc(NotFound)
	httpRouter.MethodNotAllowedHandler = http.HandlerFunc(InvalidMethod)
	httpRouter.HandleFunc(PATH, Print).Methods(http.MethodGet)
	httpRouter.HandleFunc(PATH, Create).Methods(http.MethodPost)
	httpRouter.HandleFunc(PATH, Edit).Methods(http.MethodPut)

	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", PORT),
		WriteTimeout:      time.Second * 10,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 20,
		Handler:           httpRouter,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("server error occurred '%s'\n", err.Error())
		}
	}()

	interruptAndExit(server)
}

func interruptAndExit(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	<-interruptChan

	fmt.Printf("server stopping...\n")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// #nosec G104 -- There could be an error while shutting down, but its during shutdown
	server.Shutdown(ctx)
}
