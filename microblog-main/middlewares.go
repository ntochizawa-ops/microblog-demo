package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func logRequest() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				log.Info().Dict("response", zerolog.Dict().
					Int("status", ww.Status()).
					Int("sent_bytes", ww.BytesWritten()).
					Int64("elapsed", time.Since(t1).Microseconds()),
				).Dict("request", zerolog.Dict().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("uri", r.URL.String()).
					Str("host", r.Host).
					Str("remote", r.RemoteAddr),
				).Msgf("")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
