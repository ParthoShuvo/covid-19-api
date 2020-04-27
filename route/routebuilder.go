// this package is inspired from slurpy project created by Steef de Rooi (S.deRooi@ibfd.org)

package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Builder holds all routes
type Builder struct {
	allowCROS     bool
	appName       string
	allowLogDebug bool
	router        *mux.Router
}

// Router returns the configured route
func (rb *Builder) Router() *mux.Router {
	return rb.router
}

// NewRouteBuilder constructs RouterBuilder
func NewRouteBuilder(allowCROS bool, appName string, allowLogDebug bool) *Builder {
	return &Builder{allowCROS, appName, allowLogDebug, mux.NewRouter().StrictSlash(true)}
}

// NewSubrouteBuilder constructs a subroute by a path prefix
func (rb *Builder) NewSubrouteBuilder(pathPrefix string) *Builder {
	return &Builder{
		rb.allowCROS,
		rb.appName,
		rb.allowLogDebug,
		rb.router.PathPrefix(pathPrefix).Subrouter(),
	}
}

// Add a route
func (rb *Builder) Add(name string, methods []string, path string, queries map[string]string, handler http.HandlerFunc) *mux.Route {
	return rb.add(name, methods, path, queries, handler)
}

func (rb *Builder) add(name string, methods []string, path string, queries map[string]string, handler http.HandlerFunc) *mux.Route {
	return rb.router.Name(name).Methods(methods...).Path(path).Queries(rb.toPairs(queries)...).Handler(handler)
}

func (rb *Builder) toPairs(queries map[string]string) []string {
	pairs := []string{}
	if queries == nil {
		return pairs
	}
	for name, pattern := range queries {
		pairs = append(pairs, name, pattern)
	}
	return pairs
}
