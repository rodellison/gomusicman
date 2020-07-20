package common

import (
	"fmt"
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/models"
	"strconv"
	"strings"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func ConvertDate(dateValue string) string {

	date := dateValue
	t, _ := time.Parse(layoutISO, date)
	return t.Format(layoutUS)

}

//This function calculates a min and maxDate range that will be used for Location oriented calendar filtering
//If the user provides a month in their Alexa query, use it, but if not (i.e. open ended search..) then force the results
//to filter just the current month.. as there will still be many...
func GetDatesForCalendarMinMax(requestMonth string) (string, string) {

	names := map[string]int{"January": 1, "February": 2, "March": 3, "April": 4, "May": 5, "June": 6,
		"July": 7, "August": 8, "September": 9, "October": 10, "November": 11, "December": 12}

	var requestedMonthOrd int
	if requestMonth == "" {
		//default to current month
		requestedMonthOrd = int(time.Now().Month())
	} else {
		requestedMonthOrd = names[requestMonth]
	}

	year, month, _ := time.Now().Date()
	var m int = int(month) // normally written as 'i := int(m)'

	var minDate, maxDate string
	if requestedMonthOrd < m {
		year += 1
		minDate = strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", requestedMonthOrd) + "-01"
		maxDate = strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", requestedMonthOrd) + "-31"
	} else {
		minDate = strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", requestedMonthOrd) + "-01"
		maxDate = strconv.Itoa(year) + "-" + fmt.Sprintf("%02d", requestedMonthOrd) + "-31"
	}

	return minDate, maxDate

}

func CheckDynamoForCorrectedValue(value string) string {

	strValue := cleanupKnownUserError(value)
	strValue = clients.QueryMusicManParmTable(strValue)

	return strValue

}

func cleanupKnownUserError(theValue string) string {
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
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), "the rock group ", "", 1)
	cleanedUpValue = strings.Replace(strings.ToLower(cleanedUpValue), "the music group ", "", 1)

	return cleanedUpValue

}

func ConvertStateAbbreviation(stateLocation string) string {

	stateLocation = stateLocation[0:strings.LastIndex(stateLocation, ",")]
	//Songkick uses the State abbreviation so convert it. The state is now the LAST two chars in this string..
	thisStateLoc := strings.LastIndex(stateLocation, ",")
	value := stateLocation[0:thisStateLoc] + " " + models.USC[stateLocation[thisStateLoc+2:]]
	return value

}

func UniqueEvents(eventsSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range eventsSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
