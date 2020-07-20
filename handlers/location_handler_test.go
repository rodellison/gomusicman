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

func TestFetchLocationData(t *testing.T) {

	location := make([]models.Location, 1)
	location[0].City.DisplayName = "Dayton"
	location[0].City.State.DisplayName = "OH"
	location[0].MetroArea.ID = 12345
	locationResults := models.LocationResults{Location: location}
	locationResultsPage := models.LocationResultsPage{
		Status:       "ok",
		Results:      locationResults,
		TotalEntries: 1,
	}

	var locationIDResponse models.LocationIDResponse = models.LocationIDResponse{ResultsPage: locationResultsPage}
	APIRequestLocationID = func(string) (*models.LocationIDResponse, error) {
		return &locationIDResponse, nil
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

	var locationCalendarResponse models.CalendarResponse = models.CalendarResponse{ResultsPage: calendarResultsPage}
	APIRequestLocationEventCalendar = func(string) (*models.CalendarResponse, error) {
		return &locationCalendarResponse, nil
	}

	locationCalendarEvents, _ := fetchLocationData("Dayton", "Ohio", "", "")
	assert.Contains(t, locationCalendarEvents[0], "July 1, 2020")
}

/*
func TestRealLocationQuery(t *testing.T) {

	clients.TheHTTPClient = &http.Client{}
	APIRequestLocationID = apiRequestLocationID
	APIRequestLocationEventCalendar = apiRequestArtistEventCalendar

	locationCalendarEvents, _ := fetchLocationData("New York City", "")
	assert.NotNil(t, locationCalendarEvents)

}

*/

//TODO Add a test for HandleArtistIntent
