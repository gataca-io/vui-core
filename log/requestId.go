package log

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/random"
)

const (
	HeaderXSpanId = "X-Span-Id"
	CtxTraceId    = "TRACE_ID"
	CtxSpanId     = "SPAN_ID"
)

type (
	// RequestIDConfig defines the config for RequestID middleware.
	RequestIDConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		// Generator defines a function to generate an ID.
		// Optional. Default value random.String(32).
		Generator     func() string
		TraceIdHeader string
		SpanIdHeader  string
	}
)

var (
	// DefaultRequestIDConfig is the default RequestID middleware config.
	DefaultRequestIDConfig = RequestIDConfig{
		Skipper:       middleware.DefaultSkipper,
		Generator:     generator,
		TraceIdHeader: echo.HeaderXRequestID,
		SpanIdHeader:  HeaderXSpanId,
	}
)

// RequestID returns a X-Request-ID middleware.
func RequestID() echo.MiddlewareFunc {
	return RequestIDWithConfig(DefaultRequestIDConfig)
}

// RequestID returns a X-Request-ID middleware with partial config.
func RequestIDWithHeaders(traceIdHeader string, spanIdHeader string) echo.MiddlewareFunc {
	Config := RequestIDConfig{
		Skipper:       middleware.DefaultSkipper,
		Generator:     generator,
		TraceIdHeader: traceIdHeader,
		SpanIdHeader:  spanIdHeader,
	}
	return RequestIDWithConfig(Config)
}

// RequestIDWithConfig returns a X-Request-ID middleware with config.
func RequestIDWithConfig(config RequestIDConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRequestIDConfig.Skipper
	}
	if config.Generator == nil {
		config.Generator = generator
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(config.TraceIdHeader)
			if rid == "" {
				rid = config.Generator()
			}
			res.Header().Set(config.TraceIdHeader, rid)
			c.Set(CtxTraceId, rid)
			sid := req.Header.Get(config.SpanIdHeader)
			if sid != "" {
				res.Header().Set(config.SpanIdHeader, rid)
				c.Set(CtxSpanId, sid)
			}
			return next(c)
		}
	}
}

func GetTraceId(c echo.Context) string {
	if c != nil && c.Get(CtxTraceId) != nil {
		return c.Get(CtxTraceId).(string)
	}
	return "-"
}

func GetSpanId(c echo.Context) string {
	if c != nil && c.Get(CtxSpanId) != nil {
		return c.Get(CtxSpanId).(string)
	}
	return "-"
}

func generator() string {
	return random.String(32)
}
