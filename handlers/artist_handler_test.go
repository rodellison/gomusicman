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

func TestConvertStateAbbreviation(t *testing.T) {

	location := "Virginia Beach, VA, US"
	expectedResult := "Virginia Beach, Virginia"
	//	var artistIDResponse models.ArtistIDResponse
	response := convertStateAbbreviation(location)
	assert.Equal(t, expectedResult, response)

}

func TestAPIRequestArtistID(t *testing.T) {

	requestURL := "http://api.songkick.com/api/3.0/search/artists.json?query=Iron+Maiden"
	JSONResult := "{\"resultsPage\":{\"status\":\"ok\",\"results\":{\"artist\":[{\"id\":438390,\"displayName\":\"Iron Maiden\",\"uri\":\"https://www.songkick.com/artists/438390-iron-maiden?utm_source=40512&utm_medium=partner\",\"onTourUntil\":\"2021-07-11\"}]},\"perPage\":50,\"page\":1,\"totalEntries\":1}}"

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
	//	var artistIDResponse models.ArtistIDResponse
	artistIDResponse, _ := APIRequestArtistID(requestURL)
	assert.Equal(t, artistIDResponse.ResultsPage.Results.Artist[0].DisplayName, "Iron Maiden")

}

func TestAPIRequestArtistCalendar(t *testing.T) {

	requestURL := "http://api.songkick.com/api/3.0/artists/438390/calendar.json"
	JSONResult := "{\"resultsPage\":{\"status\":\"ok\",\"results\":{\"event\":[{\"id\":39348654,\"displayName\":\"Iron Maiden with Avatar at Hallenstadion (July 13, 2020) (CANCELLED) \",\"type\":\"Concert\",\"uri\":\"http://www.songkick.com/concerts/39348654-iron-maiden-at-hallenstadion?utm_source=40512&utm_medium=partner\",\"status\":\"cancelled\",\"popularity\":0.166349,\"start\":{\"date\":\"2020-07-13\",\"datetime\":\"2020-07-13T19:00:00+0200\",\"time\":\"19:00:00\"},\"ageRestriction\":null,\"flaggedAsEnded\":false,\"venue\":{\"id\":29550,\"displayName\":\"Hallenstadion\",\"uri\":\"http://www.songkick.com/venues/29550-hallenstadion?utm_source=40512&utm_medium=partner\",\"metroArea\":{\"displayName\":\"Zürich\",\"country\":{\"displayName\":\"Switzerland\"},\"id\":104761,\"uri\":\"http://www.songkick.com/metro-areas/104761-switzerland-zurich?utm_source=40512&utm_medium=partner\"},\"lat\":47.41161,\"lng\":8.55166},\"location\":{\"city\":\"Zürich, Switzerland\"}}]},\"perPage\":50,\"page\":1,\"totalEntries\":1}}"

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
	//	var artistCalendarResponse models.CalendarResponse
	artistCalendarResponse, _ := APIRequestArtistEventCalendar(requestURL)
	assert.Equal(t, artistCalendarResponse.ResultsPage.Results.Event[0].Location.City, "Zürich, Switzerland")

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
