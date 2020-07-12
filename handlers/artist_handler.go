package handlers

import (
	"encoding/json"
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/models"
	"io/ioutil"
)

func APIRequestArtistID(urlToGet string) (models.ArtistIDResponse, error) {

	response, err := clients.GetURL(urlToGet)
	if err != nil {
		return models.ArtistIDResponse{}, nil
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var artistIDReponse models.ArtistIDResponse
		json.Unmarshal(data, &artistIDReponse)
		return artistIDReponse, nil
	}
}

func APIRequestArtistEventCalendar(urlToGet string) (models.CalendarResponse, error) {

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
