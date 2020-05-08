package cssedaily

type (
	RegionReport struct {
		CountryRegion string `json:"country/region"`
		ProvinceState string `json:"province/state"`
		CC            string `json:"cc"`
		LastUpdate    string `json:"last-update"`
		Confirmed     int    `json:"confirmed"`
		Deaths        int    `json:"deaths"`
		Recovered     int    `json:"recovered"`
	}

	DailyReport struct {
		Date           string          `json:"date"`
		TotalConfirmed int             `json:"total-confirmed"`
		TotalDeaths    int             `json:"total-deaths"`
		TotalRecovered int             `json:"total-recovered"`
		Reports        []*RegionReport `json:"reports"`
	}
	CsseDailyDatastore interface {
		ReadAllDailyReports() ([]*DailyReport, error)
		ReadDailyReport(date string) (*DailyReport, error)
	}
)
