package handlers

import (
	"bytes"
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/mocks"
	"github.com/rodellison/gomusicman/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func init() {
	clients.TheHTTPClient = &mocks.MockHTTPClient{}
}

func TestAPIRequestVenueID(t *testing.T) {

	requestURL := "http://api.songkick.com/api/3.0/search/venues.json?query=Staples+Center"
	JSONResult := "{\"resultsPage\":{\"status\":\"ok\",\"results\":{\"venue\":[{\"id\":598,\"displayName\":\"STAPLES Center\",\"uri\":\"http://www.songkick.com/venues/598-staples-center?utm_source=40512&utm_medium=partner\",\"metroArea\":{\"displayName\":\"Los Angeles\",\"country\":{\"displayName\":\"US\"},\"state\":{\"displayName\":\"CA\"},\"id\":17835,\"uri\":\"http://www.songkick.com/metro_areas/17835-us-los-angeles?utm_source=40512&utm_medium=partner\"},\"city\":{\"displayName\":\"Los Angeles\",\"country\":{\"displayName\":\"US\"},\"state\":{\"displayName\":\"CA\"},\"id\":17835}}]},\"perPage\":50,\"page\":1,\"totalEntries\":5}}"

	// build response html
	// create a new reader with that html
	mocks.GetDoHTTPFunc = func(*http.Request) (*http.Response, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		r := ioutil.NopCloser(bytes.NewReader([]byte(JSONResult)))
		return &http.Response{
			StatusCode: 200, //for this test, just using a bad return code to signify http get error
			Body:       r,
		}, nil
	}
	venueIDResponse, _ := APIRequestVenueID(requestURL)
	assert.Equal(t, venueIDResponse.ResultsPage.Results.Venue[0].City.DisplayName, "Los Angeles")

}

func TestAPIRequestVenueCalendar(t *testing.T) {

	requestURL := "http://api.songkick.com/api/3.0/venues/598/calendar.json"
	JSONResult := "{\"resultsPage\":{\"status\":\"ok\",\"results\":{\"event\":[{\"id\":39271999,\"displayName\":\"Camila Cabello with PRETTYMUCH at STAPLES Center (August 7, 2020) (CANCELLED) \",\"type\":\"Concert\",\"uri\":\"http://www.songkick.com/concerts/39271999-camila-cabello-at-staples-center?utm_source=40512&utm_medium=partner\",\"status\":\"cancelled\",\"popularity\":0.213557,\"start\":{\"date\":\"2020-08-07\",\"datetime\":null,\"time\":null},\"ageRestriction\":null,\"flaggedAsEnded\":false,\"venue\":{\"id\":598,\"displayName\":\"STAPLES Center\",\"uri\":\"http://www.songkick.com/venues/598-staples-center?utm_source=40512&utm_medium=partner\",\"metroArea\":{\"displayName\":\"Los Angeles (LA)\",\"country\":{\"displayName\":\"US\"},\"state\":{\"displayName\":\"CA\"},\"id\":17835,\"uri\":\"http://www.songkick.com/metro-areas/17835-us-los-angeles-la?utm_source=40512&utm_medium=partner\"},\"lat\":34.04301,\"lng\":-118.26736},\"location\":{\"city\":\"Los Angeles (LA), CA, US\",\"lat\":34.04301,\"lng\":-118.26736}}]},\"perPage\":50,\"page\":1,\"totalEntries\":1}}"

	// build response html
	// create a new reader with that html
	mocks.GetDoHTTPFunc = func(*http.Request) (*http.Response, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		r := ioutil.NopCloser(bytes.NewReader([]byte(JSONResult)))
		return &http.Response{
			StatusCode: 200, //for this test, just using a bad return code to signify http get error
			Body:       r,
		}, nil
	}
	venueCalendarResponse, _ := APIRequestVenueEventCalendar(requestURL)
	assert.Equal(t, venueCalendarResponse.ResultsPage.Results.Event[0].Location.City, "Los Angeles (LA), CA, US")

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
