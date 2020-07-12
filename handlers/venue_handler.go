package handlers

import (
	"encoding/json"
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/models"
	"io/ioutil"
)

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