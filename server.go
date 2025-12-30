package go_service_bootstrap_server_echo

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HTTPServerEcho struct {
	Echo *echo.Echo
	Name string
	Port int
}

func (s *HTTPServerEcho) Init() {
	s.Echo = echo.New()
	s.Echo.Use(middleware.Recover())
}

func (s *HTTPServerEcho) EnableDebug() {
	if s.Echo == nil {
		s.Init()
	}
	s.Echo.Debug = true
}

func (s *HTTPServerEcho) SetLogger(config middleware.RequestLoggerConfig) {
	if s.Echo == nil {
		s.Init()
	}
	s.Echo.Use(middleware.RequestLoggerWithConfig(config))
	s.Echo.Logger.SetPrefix(fmt.Sprintf("[HTTPServerEcho: %s]", s.Name))
}

func (s *HTTPServerEcho) Start(ctx context.Context) error {
	if s.Echo == nil {
		s.Init()
	}
	defer func() {
		recoverInfo := recover()
		if recoverInfo != nil {
			s.Echo.Logger.Fatal("recover", recoverInfo)
		}
		s.Echo.Logger.Info(
			"Server Shutdown",
			s.Echo.Shutdown(context.Background()),
		)
	}()
	var quitSignal = make(chan bool, 1)
	var startErr error
	go func() {
		defer func() {
			quitSignal <- true
			close(quitSignal)
		}()
		s.Echo.Logger.Info("Start", s.Port)
		startErr = s.Echo.Start(":" + strconv.Itoa(s.Port))
		s.Echo.Logger.Info("Start Quit", startErr)
	}()
	for {
		select {
		case <-ctx.Done():
			s.Echo.Logger.Debug("Input Context Done")
			return ctx.Err()
		case <-quitSignal:
			s.Echo.Logger.Debug("Get Quit Signal")
			return startErr
		}
	}
}

func (s *HTTPServerEcho) ServerName() string {
	return s.Name
}
