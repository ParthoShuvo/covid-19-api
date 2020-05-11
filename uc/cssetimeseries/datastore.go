package cssetimeseries

type (
	Report struct {
		CountryRegion string  `json:"country/region"`
		ProvinceState string  `json:"province/state"`
		Lat           float64 `json:"latitude"`
		Long          float64 `json:"Longitude"`
		Confirmed     int     `json:"Confirmed"`
		Deaths        int     `json:"Deaths"`
		Recovered     int     `json:"Recovered"`
	}

	Date string

	TimeSeries map[Date][]*Report

	CsseTimeSeriesData struct {
		Total      int        `json:"total"`
		Page       int        `json:"page"`
		PageSize   int        `json:"page-size"`
		TimeSeries TimeSeries `json:"time-series"`
	}

	CsseTimeSeriesDatastore interface {
		ReadCsseTimeSeries(times []string) (TimeSeries, error)
	}
)
