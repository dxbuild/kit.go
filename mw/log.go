package mw

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	TagMethod  = "method"
	TagLatency = "latency"
	TagStatus  = "status"
	TagPath    = "path"
	TagURL     = "url"
)

func NewLogger(log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			start := time.Now()
			nextErr := next(c)
			if nextErr != nil {
				c.Error(nextErr)
			}
			stop := time.Now()
			req := c.Request()
			res := c.Response()
			var evt *zerolog.Event
			switch {
			case res.Status >= 500:
				evt = log.Error()
			case res.Status >= 400:
				evt = log.Warn()
			default:
				evt = log.Info()
			}
			evt = evt.Str(TagMethod, req.Method)
			evt = evt.Str(TagPath, req.URL.Path)
			evt = evt.Str(TagURL, req.URL.String())
			evt = evt.Int(TagStatus, res.Status)
			evt = evt.Dur(TagLatency, stop.Sub(start))
			evt.Msg(http.ConnState(res.Status).String())
			return nil
		}
	}
}
