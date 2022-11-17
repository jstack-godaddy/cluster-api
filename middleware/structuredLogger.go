package middleware

import (
	"bytes"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DefaultStructuredLogger logs a gin HTTP request in JSON format. Uses the
// default logger from rs/zerolog.
func DefaultStructuredLogger() gin.HandlerFunc {
	return StructuredLogger(&log.Logger)
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// StructuredLogger logs a gin HTTP request in JSON format. Allows to set the
// logger for testing purposes.
func StructuredLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		lf, _ := os.Create("/tmp/gin.log")
		fileLogger := zerolog.New(lf).With().Logger()

		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		var logEvent *zerolog.Event
		if c.Writer.Status() >= 500 {
			logEvent = fileLogger.Error()
		} else {
			logEvent = fileLogger.Info()
		}

		ignored_paths := []string{
			"/v1/swagger/doc.json",
			"/v1/swagger/index.html",
			"/v1/swagger/swagger-ui.css",
			"/v1/swagger/favicon-32x32.png",
			"/v1/swagger/swagger-ui-bundle.js",
			"/v1/swagger/swagger-ui-standalone-preset.js",
		}

		if contains(ignored_paths, param.Path) {
			logEvent.Str("client_id", param.ClientIP).
				Str("method", param.Method).
				Int("status_code", param.StatusCode).
				Int("body_size", param.BodySize).
				Str("path", param.Path).
				Str("latency", param.Latency.String()).
				Msg(param.ErrorMessage)
		} else {
			logEvent.Str("client_id", param.ClientIP).
				Str("method", param.Method).
				Int("status_code", param.StatusCode).
				Int("body_size", param.BodySize).
				Str("path", param.Path).
				Str("latency", param.Latency.String()).
				Str("response: ", blw.body.String()).
				Msg(param.ErrorMessage)
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
