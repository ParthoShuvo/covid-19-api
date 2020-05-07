package country

import (
	"strings"

	"github.com/ParthoShuvo/fpingo/collection/list"
)

type CountryEnv struct {
	cs Countrystore
}

func NewEnv(countryStore Countrystore) *CountryEnv {
	return &CountryEnv{countryStore}
}

func (env *CountryEnv) ReadCountries(names []string) ([]*Country, error) {
	countries, err := env.ReadAllCountries()
	if len(names) == 0 {
		return countries, err
	}
	nameList := list.FromArray(names)
	predicate := func(c interface{}) bool {
		countryName := c.(*Country).Name
		finder := func(n interface{}) bool { return strings.EqualFold(countryName, n.(string)) }
		return nameList.Exists(finder)
	}

	result := []*Country{}
	for _, c := range list.FromArray(countries).Filter(predicate).ToArray() {
		result = append(result, c.(*Country))
	}
	return result, nil
}

func (env *CountryEnv) ReadAllCountries() ([]*Country, error) {
	return env.cs.ReadAllCountries()
}

func (env *CountryEnv) ReadCountryByCC(cc string) (*Country, error) {
	return env.cs.ReadCountryByCC(cc)
}

func (env *CountryEnv) ReadCountryByName(name string) (*Country, error) {
	return env.cs.ReadCountryByName(name)
}
