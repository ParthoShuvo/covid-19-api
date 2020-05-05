package country

type Country struct {
	Name                   string  `json:"name"`
	CC                     string  `json:"cc"`
	Alpha2                 string  `json:"-"`
	Alpha3                 string  `json:"-"`
	CountryCode            string  `json:"-"`
	ISO3166                string  `json:"-"`
	Region                 string  `json:"region"`
	SubRegion              string  `json:"sub-region"`
	IntermediateRegion     string  `json:"-"`
	RegionCode             string  `json:"-"`
	SubRegionCode          string  `json:"-"`
	IntermediateRegionCode string  `json:"-"`
	Latitude               float64 `json:"latitude,omitempty"`
	Longitude              float64 `json:"longitude,omitempty"`
}

type Countrystore interface {
	ReadAllCountries() ([]*Country, error)
	ReadCountryByCC(cc string) (*Country, error)
	ReadCountryByName(name string) (*Country, error)
}
