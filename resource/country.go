package resource

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/errors"
	"github.io/covid-19-api/resource/writer"
	"github.io/covid-19-api/uc/country"
)

// CountryResource defines country resources
type CountryResource struct {
	env    *country.CountryEnv
	writer writer.Writer
}

// NewCountryResource definition
func NewCountryResource(env *country.CountryEnv, w writer.Writer) *CountryResource {
	return &CountryResource{env, w}
}

// CountryFetcher provides action to fetch all countries
func (res *CountryResource) CountryFetcher() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if len(req.URL.Query()) > 0 {
			err := errors.BadRequest.New("invalid query params")
			log.Error(err)
			SendError(rw, err)
			return
		}
		countries, err := res.env.ReadAllCountries()
		if err != nil {
			log.Error(err)
			return
		}
		res.writer.Write(rw, countries)
	}
}

// CountryFetcherByCC provides action to fetch a country info by cc
func (res *CountryResource) CountryFetcherByCC() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cc := mux.Vars(req)["cc"]
		country, err := res.env.ReadCountryByCC(cc)
		if err != nil {
			log.Error(err)
			SendError(rw, err)
			return
		}
		res.writer.Write(rw, country)
	}
}

// CountryFetcherByName provides action to fetch a country info by name
func (res *CountryResource) CountryFetcherByName() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		names := req.URL.Query()["name"]
		if names == nil {
			err := errors.BadRequest.New("name query param is empty")
			log.Error(err)
			SendError(rw, err)
			return
		}
		countries, err := res.env.ReadCountries(names)
		if err != nil {
			log.Error(err)
			return
		}
		res.writer.Write(rw, countries)
	}
}
