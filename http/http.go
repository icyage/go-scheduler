package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/obsurvive/voyager/config"
	"github.com/obsurvive/voyager/http/handlers"
	midlog "github.com/obsurvive/voyager/http/middleware/log"
	"github.com/obsurvive/voyager/log"
	"github.com/obsurvive/voyager/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	echoServer *echo.Echo
	cfg        config.Provider
	sigs       chan os.Signal
)

func init() {
	cfg = config.Config()
}

func createServer() {
	echoServer = echo.New()
	echoServer.HideBanner = true
	echoServer.HidePort = true

}

func setRoutes() {
	echoServer.GET("/healthz", handlers.Healthz, middleware.BasicAuth(utils.InternalUser))
	echoServer.GET("/metrics", echo.WrapHandler(promhttp.Handler()), middleware.BasicAuth(utils.InternalUser))

	echoServer.POST("/schedule/new", handlers.ScheduleCreate)
	echoServer.GET("/schedule/:id", handlers.ScheduleGet)
	echoServer.DELETE("/schedule/:id", handlers.ScheduleDelete)
	echoServer.GET("/schedule/status", handlers.ScheduleStatus)
	echoServer.PUT("/schedule/:id", handlers.ScheduleTrigger)
}

func setMiddleware() {
	echoServer.Use(midlog.Logger())
	echoServer.Use(middleware.RequestID())
	echoServer.Use(middleware.Gzip())
	echoServer.Use(middleware.Recover())
}

func Serve() {
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createServer()
	setMiddleware()
	setRoutes()

	// Start server
	go func() {
		log.WithFields(log.Fields{
			"EventName": "server_start",
			"PID":       os.Getpid(),
			"PORT":      cfg.GetString("bind_address"),
		}).Infof("Running with PID: %i", os.Getpid())
		if err := echoServer.Start(cfg.GetString("bind_address")); err != nil {
			log.WithFields(log.Fields{
				"EventName": "server_start",
				"Error":     err,
			}).Fatalf("Shutting down the server: %s", err)
		}
	}()

	signalReceived := <-sigs
	log.WithFields(log.Fields{
		"EventName": "server_stop",
		"Signal":    signalReceived,
	}).Infof("Shutting down server: %s", signalReceived)
	if err := echoServer.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"EventName": "server_start",
			"Error":     err,
		}).Fatalf("Shutting down server: %s", err)
	}
}
