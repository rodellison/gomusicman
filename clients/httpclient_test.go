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
	urlRequest, _ := ConstructURLRequest("ArtistQuery", "Iron Maiden")
	assert.Contains(t, urlRequest, urlExpected)
	assert.Contains(t, urlRequest, queryParmExpected)

}

func TestConstructURLRequestArtistCalendar(t *testing.T) {

	urlExpected := "http://api.songkick.com/api/3.0/artists/438390/calendar.json?apikey="
	urlRequest, _ := ConstructURLRequest("ArtistCalendar", "438390")
	assert.Contains(t, urlRequest, urlExpected)

}

func TestConstructURLRequestVenueID(t *testing.T) {

	urlExpected := "http://api.songkick.com/api/3.0/search/venues.json?apikey="
	queryParmExpected := "&query=Staples+Center"
	urlRequest, _ := ConstructURLRequest("VenueQuery", "Staples Center")
	assert.Contains(t, urlRequest, urlExpected)
	assert.Contains(t, urlRequest, queryParmExpected)

}

func TestConstructURLRequestVenueCalendar(t *testing.T) {

	urlExpected := "http://api.songkick.com/api/3.0/venues/123456/calendar.json?apikey="
	urlRequest, _ := ConstructURLRequest("VenueCalendar", "123456")
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
