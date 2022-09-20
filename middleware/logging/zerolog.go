package logging

import (
	"net"

	authmiddleware "github.com/jakob-moeller-cloud/octi-sync-server/middleware/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func RequestLogging(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogError:         true,
		LogLatency:       true,
		LogMethod:        true,
		LogContentLength: true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogResponseSize:  true,
		LogUserAgent:     true,
		LogValuesFunc: func(context echo.Context, values middleware.RequestLoggerValues) error {
			logger.Debug().
				Str("Method", values.Method).
				Str("URI", values.URI).
				Int("status", values.Status).
				Str("content-length", values.ContentLength).
				Str("x-request-id", values.RequestID).
				Str("x-device-id", context.Request().Header.Get(authmiddleware.DeviceIDHeader)).
				Int64("response-size", values.ResponseSize).
				Str("user-agent", values.UserAgent).
				Err(values.Error).
				Dur("latency", values.Latency).
				IPAddr("remote-ip", net.ParseIP(values.RemoteIP)).
				Msg("request")
			return nil
		},
	})
}
