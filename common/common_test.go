package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertStateAbbreviation(t *testing.T) {

	location := "Virginia Beach, VA, US"
	expectedResult := "Virginia Beach Virginia"
	//	var artistIDResponse models.ArtistIDResponse
	response := ConvertStateAbbreviation(location)
	assert.Equal(t, expectedResult, response)

}

func TestGetDatesForCalendarMinMax(t *testing.T) {

	month := ""
	minDate, maxDate := GetDatesForCalendarMinMax(month)

	assert.Equal(t, "2020-07-01", minDate)
	assert.Equal(t, "2020-07-31", maxDate)

}
