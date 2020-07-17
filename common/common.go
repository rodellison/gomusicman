package common

import (
	"github.com/rodellison/gomusicman/clients"
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
