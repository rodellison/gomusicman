package handlers

import (
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/mocks"
	"github.com/rodellison/gomusicman/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	clients.TheHTTPClient = &mocks.MockHTTPClient{}
}

func TestFetchVenueData(t *testing.T) {

	venue := make([]models.Venue, 1)
	venue[0].ID = 12345
	venue[0].DisplayName = "Staples Center"
	venueResults := models.VenueResults{Venue: venue}
	venueResultsPage := models.VenueResultsPage{
		Status:       "ok",
		Results:      venueResults,
		TotalEntries: 1,
	}

	var venueIDResponse models.VenueIDResponse = models.VenueIDResponse{ResultsPage: venueResultsPage}
	APIRequestVenueID = func(string) (*models.VenueIDResponse, error) {
		return &venueIDResponse, nil
	}

	calendarEvents := make([]models.CalendarEvents, 1)
	calendarEvents[0].Status = "ok"
	calendarEvents[0].DisplayName = "Some great artist at Dayton Hara Arena (2020-07-01)"
	calendarEvents[0].Start.Date = "2020-07-01"
	calendarEvents[0].Venue.DisplayName = "Dayton Hara Arena"
	calendarEvents[0].Location.City = "Dayton, OH"

	calendarResults := models.CalendarResults{
		Event: calendarEvents,
	}
	calendarResultsPage := models.CalendarResultsPage{
		Status:       "success",
		Results:      calendarResults,
		Page:         1,
		PerPage:      50,
		TotalEntries: 1,
	}

	var venueCalendarResponse models.CalendarResponse = models.CalendarResponse{ResultsPage: calendarResultsPage}
	APIRequestVenueEventCalendar = func(string) (*models.CalendarResponse, error) {
		return &venueCalendarResponse, nil
	}

	venueCalendarEvents, _ := fetchVenueData("Staples Center", "July")
	assert.Contains(t, venueCalendarEvents[0], "Some great artist")
	assert.Contains(t, venueCalendarEvents[0], "July 1, 2020")
}

//TODO Add a test for HandleArtistIntent
