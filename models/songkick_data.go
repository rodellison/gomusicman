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
