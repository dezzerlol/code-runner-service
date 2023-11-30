package http

import (
	"code-runner-service/internal/mqueue"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	mq     *mqueue.MQueue
}

func New(logger *zap.SugaredLogger, mq *mqueue.MQueue) *Server {
	return &Server{
		logger: logger,
		mq:     mq,
	}
}

func (s *Server) Run(host string, port string) {
	e := echo.New()

	s.LoadRoutes(e)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			s.logger.Infow("request",
				"URI", v.URI,
				"status", v.Status,
			)

			return nil
		},
	}))

	addr := fmt.Sprintf("%s:%s", host, port)

	s.logger.Fatal(e.Start(addr))
}
