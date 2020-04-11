package resource

import (
	"net/http"

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

// CountryFetcher provies action to fetch all countries
func (res *CountryResource) CountryFetcher() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		countries := res.da.GetAll().([]db.Country)
		res.writer.Write(rw, countries)
	}
}
