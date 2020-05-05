package country

type CountryEnv struct {
	cs Countrystore
}

func NewEnv(countryStore Countrystore) *CountryEnv {
	return &CountryEnv{countryStore}
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
