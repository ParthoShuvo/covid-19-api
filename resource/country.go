package resource

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.io/covid-19-api/resource/writer"

	"github.io/covid-19-api/db"
)

// CountryResource defines country resources
type CountryResource struct {
	da     db.DataAccessor
	writer writer.Writer
}

// NewCountryResource definition
func NewCountryResource(dataAccessor db.DataAccessor, w writer.Writer) *CountryResource {
	return &CountryResource{dataAccessor, w}
}

// CountryFetcher provides action to fetch all countries
func (res *CountryResource) CountryFetcher() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		countries := res.da.GetAll().([]db.Country)
		res.writer.Write(rw, countries)
	}
}

// CountryFetcherByCC provides action to fetch a country info by cc
func (res *CountryResource) CountryFetcherByCC() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		cc := mux.Vars(req)["cc"]
		pred := func(country db.Country) bool {
			return strings.EqualFold(country.LatLong.CC, cc)
		}
		if country := res.findCountry(pred); country != nil {
			res.writer.Write(rw, country)
			return
		}
		log.Printf("no country found by cc=%s", cc)
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, "Not found")
	}
}

// CountryFetcherByName provides action to fetch a country info by name
func (res *CountryResource) CountryFetcherByName() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		name := mux.Vars(req)["name"]
		pred := func(country db.Country) bool {
			return strings.EqualFold(country.Info.Name, name)
		}
		if country := res.findCountry(pred); country != nil {
			res.writer.Write(rw, country)
			return
		}
		log.Printf("no country found by name=%s", name)
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, "Not found")
	}
}

func (res *CountryResource) findCountry(predicate func(country db.Country) bool) *db.Country {
	countries := res.da.GetAll().([]db.Country)
	for _, country := range countries {
		if predicate(country) {
			return &country
		}
	}
	return nil
}
