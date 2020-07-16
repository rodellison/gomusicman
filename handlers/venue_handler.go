package handlers

import (
	"encoding/json"
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/models"
	"io/ioutil"
	"strings"
)

const VENUE_NAME_SLOT = "venue"
const VENUE_MONTH_SLOT = "month"

func APIRequestVenueID(urlToGet string) (models.VenueIDResponse, error) {

	response, err := clients.GetURL(urlToGet)
	if err != nil {
		return models.VenueIDResponse{}, nil
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var venueIDReponse models.VenueIDResponse
		json.Unmarshal(data, &venueIDReponse)
		return venueIDReponse, nil
	}
}

func APIRequestVenueEventCalendar(urlToGet string) (models.CalendarResponse, error) {

	response, err := clients.GetURL(urlToGet)
	if err != nil {
		return models.CalendarResponse{}, nil
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var calendarReponse models.CalendarResponse
		json.Unmarshal(data, &calendarReponse)
		return calendarReponse, nil
	}
}

func cleanupKnownUserErrorForVenues(theValue string) string {
	cleanedUpValue := theValue

	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), "u. s.", "US", 1) //e.g. U. S. Bank Arena should be US Bank Arena
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), "a. t. and t ", "AT&T", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), "b. b. and t.", "BB&T", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " marina", " Arena", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " farina", " Arena", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), "amplitheater", "Amphitheater", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " today", "", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " tonight", "", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " tomorrow", "", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " this week", "", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " next week", "", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " this month", "", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), " next month", "", 1)

	return cleanedUpValue

}
