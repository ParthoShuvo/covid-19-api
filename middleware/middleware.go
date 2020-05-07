package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.io/covid-19-api/cfg"
	"github.io/covid-19-api/errors"
	"github.io/covid-19-api/resource"

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

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				wrappedErr := errors.InternalServerError.New(fmt.Sprint(err))
				log.Error(wrappedErr)
				resource.SendError(rw, wrappedErr)
			}
		}()
		next.ServeHTTP(rw, r)
	})
}
