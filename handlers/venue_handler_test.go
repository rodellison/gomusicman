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
	var venueIDResponse models.VenueIDResponse
	venueIDResponse, _ = APIRequestVenueID(requestURL)
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
	var venueCalendarResponse models.CalendarResponse
	venueCalendarResponse, _ = APIRequestVenueEventCalendar(requestURL)
	assert.Equal(t, venueCalendarResponse.ResultsPage.Results.Event[0].Location.City, "Los Angeles (LA), CA, US")

}
