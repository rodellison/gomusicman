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

func TestFetchArtistData(t *testing.T) {

	artist := make([]models.Artist, 1)
	artist[0].ID = 12345
	artist[0].DisplayName = "Iron Maiden"
	artistResults := models.ArtistIDResults{Artist: artist}
	artistResultsPage := models.ArtistIDResultsPage{
		Status:       "ok",
		Results:      artistResults,
		TotalEntries: 1,
	}

	var artistIDResponse models.ArtistIDResponse = models.ArtistIDResponse{ResultsPage: artistResultsPage}
	APIRequestArtistID = func(string) (*models.ArtistIDResponse, error) {
		return &artistIDResponse, nil
	}

	calendarEvents := make([]models.CalendarEvents, 1)
	calendarEvents[0].Status = "ok"
	calendarEvents[0].DisplayName = "Iron Maiden with special guests at Dayton Hara Arena (2020-07-01)"
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

	var artistCalendarResponse models.CalendarResponse = models.CalendarResponse{ResultsPage: calendarResultsPage}
	APIRequestArtistEventCalendar = func(string) (*models.CalendarResponse, error) {
		return &artistCalendarResponse, nil
	}

	artistCalendarEvents, _ := fetchArtistData("Iron Maiden", "July")
	assert.Contains(t, artistCalendarEvents[0], "July 1, 2020")
}

//TODO Add a test for HandleArtistIntent
