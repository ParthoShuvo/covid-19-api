package middleware

import (
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.io/covid-19-api/cfg"

	log "github.com/sirupsen/logrus"
)

// PerformanceLogMiddleware logs route action processing performance
func PerformanceLogMiddleware(next http.Handler, allowLogDebug bool) http.Handler {
	if allowLogDebug {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(rw, r)
			log.Infof("%s %s %s", r.Method, r.RequestURI, time.Since(start))
		})
	}
	return next
}

// CORSMiddleware handles cors
func CORSMiddleware(next http.Handler, allowCORS bool, corsDef *cfg.CORSDef) http.Handler {
	if allowCORS {
		c := cors.New(cors.Options{
			AllowedOrigins:   corsDef.AllowedOrigins,
			AllowCredentials: corsDef.AllowCredentials,
			AllowedMethods:   corsDef.AllowedMethods,
			Debug:            corsDef.Debug,
		})
		return c.Handler(next)
	}
	return next
}
