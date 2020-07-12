package models

//Very useful json to struct utility: https://mholt.github.io/json-to-go/
//Note: have manually excluded a number of unneeded extra internal struct and property data returned by SongKick
//to minimize data size that may be carried forward in session attributes
type CalendarResponse struct {
	ResultsPage struct {
		Status  string `json:"status"`
		Results struct {
			Event []struct {
				ID          int     `json:"id"`
				DisplayName string  `json:"displayName"`
				Status      string  `json:"status"`
				Start       struct {
					Date     string      `json:"date"`
				} `json:"start"`
				Venue          struct {
					ID          int    `json:"id"`
					DisplayName string `json:"displayName"`
				} `json:"venue"`
				Location struct {
					City string  `json:"city"`
				} `json:"location"`
			} `json:"event"`
		} `json:"results"`
		TotalEntries int `json:"totalEntries"`
	} `json:"resultsPage"`
}

type ArtistIDResponse struct {
	ResultsPage struct {
		Status  string `json:"status"`
		Results struct {
			Artist []struct {
				ID          int    `json:"id"`
				DisplayName string `json:"displayName"`
			} `json:"artist"`
		} `json:"results"`
		TotalEntries int `json:"totalEntries"`
	} `json:"resultsPage"`
}

type VenueIDResponse struct {
	ResultsPage struct {
		Status  string `json:"status"`
		Results struct {
			Venue []struct {
				ID          int    `json:"id"`
				DisplayName string `json:"displayName"`
				City        struct {
					DisplayName string `json:"displayName"`
					Country     struct {
						DisplayName string `json:"displayName"`
					} `json:"country"`
					State struct {
						DisplayName string `json:"displayName"`
					} `json:"state"`
				} `json:"city,omitempty"`
			} `json:"venue"`
		} `json:"results"`
		TotalEntries int `json:"totalEntries"`
	} `json:"resultsPage"`
}
