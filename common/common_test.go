package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertDate(t *testing.T) {
	result := ConvertDate(layoutISO)
	assert.Equal(t, result, layoutUS)
}


func TestUniqueEvents(t *testing.T) {
	duplicateStringsInArray := []string{
		"This is a test",
		"This is a test",
	}

	resultStringArray := UniqueEvents(duplicateStringsInArray)
	assert.Len(t, resultStringArray, 1)
	assert.Equal(t, resultStringArray[0], "This is a test")

}

func TestConvertStateAbbreviation(t *testing.T) {

	location := "Virginia Beach, VA, US"
	expectedResult := "Virginia Beach Virginia"
	//	var artistIDResponse models.ArtistIDResponse
	response := ConvertStateAbbreviation(location)
	assert.Equal(t, expectedResult, response)

}

func TestGetDatesForCalendarMinMax(t *testing.T) {

	month := "December"
	minDate, maxDate := GetDatesForCalendarMinMax(month)

	assert.Equal(t, "2020-12-01", minDate)
	assert.Equal(t, "2020-12-31", maxDate)

}
