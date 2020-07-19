package clients

import (
	"bytes"
	"github.com/rodellison/gomusicman/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func init() {
	TheHTTPClient = &mocks.MockHTTPClient{}
}

func TestConstructURLRequestArtistID(t *testing.T) {

	urlExpected := "http://api.songkick.com/api/3.0/search/artists.json?apikey="
	queryParmExpected := "&query=Iron+Maiden"
	urlRequest, _ := ConstructURLRequest("ArtistQuery", "Iron Maiden", "", "")
	assert.Contains(t, urlRequest, urlExpected)
	assert.Contains(t, urlRequest, queryParmExpected)

}

func TestConstructURLRequestArtistCalendar(t *testing.T) {

	urlExpected := "http://api.songkick.com/api/3.0/artists/438390/calendar.json?apikey="
	urlRequest, _ := ConstructURLRequest("ArtistCalendar", "438390", "", "")
	assert.Contains(t, urlRequest, urlExpected)

}

func TestConstructURLRequestVenueID(t *testing.T) {

	urlExpected := "http://api.songkick.com/api/3.0/search/venues.json?apikey="
	queryParmExpected := "&query=Staples+Center"
	urlRequest, _ := ConstructURLRequest("VenueQuery", "Staples Center", "", "")
	assert.Contains(t, urlRequest, urlExpected)
	assert.Contains(t, urlRequest, queryParmExpected)

}

func TestConstructURLRequestVenueCalendar(t *testing.T) {

	urlExpected := "http://api.songkick.com/api/3.0/venues/123456/calendar.json?apikey="
	urlRequest, _ := ConstructURLRequest("VenueCalendar", "123456", "", "")
	assert.Contains(t, urlRequest, urlExpected)
}

func TestGetURLReturnsGood(t *testing.T) {

	testGoodJSON := "{\"resultsPage\":{\"status\":\"ok\",\"results\":{\"artist\":[{\"id\":438390,\"displayName\":\"Iron Maiden\",\"uri\":\"https://www.songkick.com/artists/438390-iron-maiden?utm_source=40512&utm_medium=partner\",\"onTourUntil\":\"2021-07-11\"}]},\"perPage\":50,\"page\":1,\"totalEntries\":1}}"

	// build response html
	// create a new reader with that html
	mocks.GetDoHTTPFunc = func(*http.Request) (*http.Response, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		r := ioutil.NopCloser(bytes.NewReader([]byte(testGoodJSON)))
		return &http.Response{
			StatusCode: 200, //for this test, just using a bad return code to signify http get error
			Body:       r,
		}, nil
	}

	response, _ := GetURL("https://justATest.com")
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, testGoodJSON, bodyString)

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
	artistCalendarResponse, _ := APIRequestEventCalendar(requestURL)
	assert.Equal(t, artistCalendarResponse.ResultsPage.Results.Event[0].Location.City, "Zürich, Switzerland")

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
	venueCalendarResponse, _ := APIRequestEventCalendar(requestURL)
	assert.Equal(t, venueCalendarResponse.ResultsPage.Results.Event[0].Location.City, "Los Angeles (LA), CA, US")

}

func TestAPIRequestLocationID(t *testing.T) {

	requestURL := "http://api.songkick.com/api/3.0/search/locations.json?query=Fort Lauderdale"
	JSONResult := "{\"resultsPage\":{\"status\":\"ok\",\"results\":{\"location\":[{\"city\":{\"lat\":26.1358,\"lng\":-80.1418,\"country\":{\"displayName\":\"US\"},\"state\":{\"displayName\":\"FL\"},\"displayName\":\"Fort Lauderdale\"},\"metroArea\":{\"lat\":26.1358,\"lng\":-80.1418,\"country\":{\"displayName\":\"US\"},\"uri\":\"http://www.songkick.com/metro_areas/19511-us-fort-lauderdale?utm_source=40512&utm_medium=partner\",\"state\":{\"displayName\":\"FL\"},\"displayName\":\"Fort Lauderdale\",\"id\":19511}}]},\"perPage\":50,\"page\":1,\"totalEntries\":1}}"

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
	locationIDResponse, _ := APIRequestLocationID(requestURL)
	assert.Equal(t, locationIDResponse.ResultsPage.Results.Location[0].MetroArea.ID, 19511)

}

func TestAPIRequestLocationCalendar(t *testing.T) {

	requestURL := "http://api.songkick.com/api/3.0/metro_areas/19511/calendar.json"
	JSONResult := "{\"resultsPage\":{\"status\":\"ok\",\"results\":{\"event\":[{\"id\":39644230,\"displayName\":\"33 Years at Tim Finnegan's Irish Pub (July 18, 2020)\",\"type\":\"Concert\",\"uri\":\"https://www.songkick.com/concerts/39644230-33-years-at-tim-finnegans-irish-pub?utm_source=40512&utm_medium=partner\",\"status\":\"ok\",\"popularity\":0.00001,\"start\":{\"date\":\"2020-07-18\",\"datetime\":\"2020-07-18T21:00:00-0400\",\"time\":\"21:00:00\"},\"performance\":[{\"id\":74981907,\"displayName\":\"33 Years\",\"billing\":\"headline\",\"billingIndex\":1,\"artist\":{\"id\":6737684,\"displayName\":\"33 Years\",\"uri\":\"https://www.songkick.com/artists/6737684-33-years?utm_source=40512&utm_medium=partner\",\"identifier\":[]}}],\"ageRestriction\":null,\"flaggedAsEnded\":false,\"venue\":{\"id\":777956,\"displayName\":\"Tim Finnegan's Irish Pub\",\"uri\":\"https://www.songkick.com/venues/777956-tim-finnegans-irish-pub?utm_source=40512&utm_medium=partner\",\"metroArea\":{\"displayName\":\"Fort Lauderdale\",\"country\":{\"displayName\":\"US\"},\"state\":{\"displayName\":\"FL\"},\"id\":19511,\"uri\":\"https://www.songkick.com/metro-areas/19511-us-fort-lauderdale?utm_source=40512&utm_medium=partner\"},\"lat\":26.42805,\"lng\":-80.07249},\"location\":{\"city\":\"Delray Beach, FL, US\",\"lat\":26.42805,\"lng\":-80.07249}}]},\"perPage\":50,\"page\":1,\"totalEntries\":217}}"

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
	locationCalendarResponse, _ := APIRequestEventCalendar(requestURL)
	assert.Equal(t, locationCalendarResponse.ResultsPage.Results.Event[0].Location.City, "Delray Beach, FL, US")

}
