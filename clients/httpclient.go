package clients

import (
	"bytes"
	"encoding/json"
	"github.com/rodellison/gomusicman/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//An HTTPClient interface declared to allow for easier Mock testing
//Basically, ensure the custom interface has definitions for the functions that need to be mocked (so
//as to not make 'real' requests
//Establish a Variable that is of the interface's type that can be used to hold the 'real' client (when
//not running tests, as well as be a variable we can set during 'test' time
//And setup an init() function that sets the variable to the 'real' interface as it's default when not testing

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	//IF we're running a test, we'll swap this variable's value to use a mock instead, but when not
	//testing, the value will be preset to ensure that it uses the 'real' httpClient interface
	TheHTTPClient HTTPClient
)

func init() {

	TheHTTPClient = &http.Client{}
}

//func GetURL fetches raw HTML data from the input url.. essentially a screen-scrape
func GetURL(url string) (*http.Response, error) {
	//Empty body for now
	jsonBytes, err := json.Marshal("")
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	//Client := &http.Client{}  using the variable/interface above to facilitate easier mock testing later
	return TheHTTPClient.Do(request)
}

//RequestFeed handles fetching external HTML site data and Unmarhalling to a struct that can be used later
//within the respective handler functions
func ConstructURLRequest(mode string, content string, minDate string, maxDate string) (string, error) {

	//ArtistQuery, content=<string containing artist name>
	//ArtistCalendar content =<artistID string>
	//VenueQuery, content=<string containing venue name>
	//VenueCalendar content =<venueID string>

	//Iron Maiden
	//http://api.songkick.com/api/3.0/artists/438390/calendar.json?apikey=SyjXrvJQG067CHtg

	var urlString string

	switch mode {
	case "ArtistQuery":
		urlString = "http://api.songkick.com/api/3.0/search/artists.json"
	case "ArtistCalendar":
		urlString = strings.Replace("http://api.songkick.com/api/3.0/artists/<artistID>/calendar.json", "<artistID>", content, 1)
	case "VenueQuery":
		urlString = "http://api.songkick.com/api/3.0/search/venues.json"
	case "VenueCalendar":
		urlString = strings.Replace("http://api.songkick.com/api/3.0/venues/<venueID>/calendar.json", "<venueID>", content, 1)
	case "LocationQuery":
		urlString = "http://api.songkick.com/api/3.0/search/locations.json"
	case "LocationCalendar":
		urlString = strings.Replace("http://api.songkick.com/api/3.0/metro_areas/<locationID>/calendar.json", "<locationID>", content, 1)
	}

	endpoint, _ := url.Parse(urlString)
	queryParams := endpoint.Query()

	//Always need to pass query parm 'apikey'
	queryParams.Set("apikey", os.Getenv("APIKey"))

	switch mode {
	case "ArtistQuery":
		queryParams.Set("query", content)
	case "VenueQuery":
		queryParams.Set("query", content)
	case "LocationQuery":
		queryParams.Set("query", content)
	case "LocationCalendar":
		queryParams.Set("min_date", minDate)
		queryParams.Set("max_date", maxDate)
	}

	endpoint.RawQuery = queryParams.Encode()
	return endpoint.String(), nil

}

func APIRequestArtistID(urlToGet string) (*models.ArtistIDResponse, error) {

	response, err := GetURL(urlToGet)
	if err != nil {
		return &models.ArtistIDResponse{}, nil
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var artistIDReponse models.ArtistIDResponse
		json.Unmarshal(data, &artistIDReponse)
		return &artistIDReponse, nil
	}
}

func APIRequestEventCalendar(urlToGet string) (*models.CalendarResponse, error) {

	response, err := GetURL(urlToGet)
	if err != nil {
		return &models.CalendarResponse{}, nil
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var calendarReponse models.CalendarResponse
		json.Unmarshal(data, &calendarReponse)
		return &calendarReponse, nil
	}
}

func APIRequestVenueID(urlToGet string) (*models.VenueIDResponse, error) {

	response, err := GetURL(urlToGet)
	if err != nil {
		return &models.VenueIDResponse{}, nil
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var venueIDReponse models.VenueIDResponse
		json.Unmarshal(data, &venueIDReponse)
		return &venueIDReponse, nil
	}
}

func APIRequestLocationID(urlToGet string) (*models.LocationIDResponse, error) {

	response, err := GetURL(urlToGet)
	if err != nil {
		return &models.LocationIDResponse{}, nil
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var locationIDReponse models.LocationIDResponse
		json.Unmarshal(data, &locationIDReponse)
		return &locationIDReponse, nil
	}
}
