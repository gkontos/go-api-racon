package controller

/**
https://github.com/quycao/chizerolog/blob/master/logger.go
**/

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type StructuredLogger struct {
	Logger *zerolog.Logger
}

// StructuredLoggerEntry is wrapper for zerolog.Event
type StructuredLoggerEntry struct {
	Logger *zerolog.Event
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: l.Logger.Info()}

	if rec := recover(); rec != nil {
		entry = &StructuredLoggerEntry{Logger: l.Logger.Error()}
	}

	logFields := map[string]interface{}{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.Fields(logFields)

	return entry
}

// Write is method that was call when server response the request
func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.Fields(map[string]interface{}{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Msg("request complete")
}

// Panic is method that was call when server have panic with request
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.Fields(map[string]interface{}{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
	l.Logger.Msg("request failed")
}
