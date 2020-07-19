package models

//Very useful json to struct utility: https://mholt.github.io/json-to-go/
//Note: have manually excluded a number of unneeded extra internal struct and property data returned by SongKick
//to minimize data size

//-----------------Event Calendar Structure

type CalendarResponse struct {
	ResultsPage CalendarResultsPage `json:"resultsPage"`
}

type CalendarResultsPage struct {
	Status       string          `json:"status"`
	Results      CalendarResults `json:"results"`
	PerPage      int             `json:"perPage"`
	Page         int             `json:"page"`
	TotalEntries int             `json:"totalEntries"`
}

type CalendarResults struct {
	Event []CalendarEvents `json:"event"`
}

type CalendarEvents struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`
	Start       struct {
		Date string `json:"date"`
	} `json:"start"`
	Venue struct {
		ID          int    `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"venue"`
	Location struct {
		City string `json:"city"`
	} `json:"location"`
}

//-----------------Artist ID Structure

type ArtistIDResponse struct {
	ResultsPage ArtistIDResultsPage `json:"resultsPage"`
}

type ArtistIDResultsPage struct {
	Status       string          `json:"status"`
	Results      ArtistIDResults `json:"results"`
	TotalEntries int             `json:"totalEntries"`
}

type ArtistIDResults struct {
	Artist []Artist `json:"artist"`
}

type Artist struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
}

//-----------------Venue ID Structure

type VenueIDResponse struct {
	ResultsPage VenueResultsPage `json:"resultsPage"`
}

type VenueResultsPage struct {
	Status       string       `json:"status"`
	Results      VenueResults `json:"results"`
	TotalEntries int          `json:"totalEntries"`
}

type VenueResults struct {
	Venue []Venue `json:"venue"`
}

type Venue struct {
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
}

//-----------------Location ID Structure

type LocationIDResponse struct {
	ResultsPage LocationResultsPage `json:"resultsPage"`
}

type LocationResultsPage struct {
	Status       string          `json:"status"`
	Results      LocationResults `json:"results"`
	PerPage      int             `json:"perPage"`
	Page         int             `json:"page"`
	TotalEntries int             `json:"totalEntries"`
}

type LocationResults struct {
	Location []Location `json:"location"`
}

type Location struct {
	City struct {
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
		Country struct {
			DisplayName string `json:"displayName"`
		} `json:"country"`
		State struct {
			DisplayName string `json:"displayName"`
		} `json:"state"`
		DisplayName string `json:"displayName"`
	} `json:"city,omitempty"`
	MetroArea struct {
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
		Country struct {
			DisplayName string `json:"displayName"`
		} `json:"country"`
		URI   string `json:"uri"`
		State struct {
			DisplayName string `json:"displayName"`
		} `json:"state"`
		DisplayName string `json:"displayName"`
		ID          int    `json:"id"`
	} `json:"metroArea,omitempty"`
}
