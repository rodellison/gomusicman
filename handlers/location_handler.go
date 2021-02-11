package handlers

import (
	"fmt"
	"github.com/rodellison/gomusicman/alexa"
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/common"
	"github.com/rodellison/gomusicman/models"
	"os"
	"strconv"
	"strings"
)

var (
	//For mocking/testing overrides
	APIRequestLocationID            func(string) (*models.LocationIDResponse, error)
	APIRequestLocationEventCalendar func(string) (*models.CalendarResponse, error)
	LOCATION_ID                     string
)

func init() {
	APIRequestLocationID = clients.APIRequestLocationID
	APIRequestLocationEventCalendar = clients.APIRequestEventCalendar
}

const (
	LOCATION_CITY_SLOT    = "city"
	LOCATION_STATE_SLOT   = "state"
	LOCATION_REGION_SLOT   = "region"
	LOCATION_COUNTRY_SLOT = "country"
	LOCATION_MONTH_SLOT   = "month"
	LOCATION_INTENT       = "LocationIntent"
)

type LocationData struct {
	ID        string
	Name      string
	Eventdata []string
}

func fetchLocationData(city, state, country, month string) ([]string, error) {

	thisMonth := strings.Title(month)
	LOCATION_ID = ""

	urlToFetch, err := clients.ConstructURLRequest("LocationQuery", city, "", "")
	if err != nil {
		return nil, err
	}
	//Make an API call to Songkick to get the Metro_area ID for this artist
	locationIDResponse, err := APIRequestLocationID(urlToFetch)
	if err != nil {
		return nil, err
	}

	if locationIDResponse.ResultsPage.TotalEntries == 0 {
		//This artist wasn't found, so return immediately..
		return nil, nil
	}

	//To get the 'correct' LocationID, we have to scan through the returned list either using the state, or country
	if state != "" {
		stateAbbr := models.USC_ABBREV[state]
		for _, item := range locationIDResponse.ResultsPage.Results.Location {
			if item.City.State.DisplayName == stateAbbr {
				LOCATION_ID = strconv.Itoa(item.MetroArea.ID)
				break
			}
		}
	}
	if country != "" {
		for _, item := range locationIDResponse.ResultsPage.Results.Location {
			if item.City.Country.DisplayName == country {
				LOCATION_ID = strconv.Itoa(item.MetroArea.ID)
				break
			}
		}
	}

	if LOCATION_ID == "" {
		fmt.Println("Could NOT find the right LOCATION ID based on input parms")
		return nil, err
	}

	var locationResults []models.CalendarEvents

	//When getting the Location oriented Calendar events.. there can be TONS!!!.. so we have to limit these to
	//a month...
	minDate, maxDate := common.GetDatesForCalendarMinMax(thisMonth)
	urlToFetch, err = clients.ConstructURLRequest("LocationCalendar", LOCATION_ID, minDate, maxDate)
	if err != nil {
		return nil, err
	}
	fmt.Println("URL being fetched: ", urlToFetch)
	//Make an API call to Songkick to get the Metro/Location's Event Calendar
	locationCalendarResponse, err := APIRequestLocationEventCalendar(urlToFetch)
	if err != nil {
		return nil, err
	} else {
		locationResults = append(locationResults, locationCalendarResponse.ResultsPage.Results.Event...)
	}

	counter := 0
	var itemsToSave []string
	var displayEventString string

	for _, item := range locationResults {
		//If the user passed a Month as part of their request.. then filter out just those events..
		//The end result may be that no events are included.
		dateString := " on " + common.ConvertDate(item.Start.Date)
		thisLocation := item.Location.City
		if strings.Contains(thisLocation, ", US") {
			//Songkick uses the State abbreviation so convert it. The state is now the LAST two chars in this string..
			thisLocation = common.ConvertStateAbbreviation(thisLocation)
		}
		if strings.Contains(item.DisplayName, " at ") {
			displayLocAt := strings.Index(item.DisplayName, " at ")
			displayEventString = item.DisplayName[0:displayLocAt] + " at " + item.Venue.DisplayName + dateString + " in " + thisLocation
		} else {
			displayEventString = item.DisplayName + " at " + item.Venue.DisplayName + dateString + " in " + thisLocation
		}
		if strings.Contains(item.DisplayName, "CANCELLED") {
			displayEventString += ", is CANCELLED."
		}

		itemsToSave = append(itemsToSave, displayEventString)
		counter += 1

	}

	return common.UniqueEvents(itemsToSave), nil

}

//Parameters passed allow this function to accommodate both the initial request, as well as subsequent
//requests as a result of the user saying 'yes' for more data
func HandleLocationIntent(request alexa.Request, resumingPrior bool, sessionData models.SessionData) alexa.Response {

	var eventData []string
	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var cardTextContent string
	var strLocationCity string
	var strLocationState string
	var strLocationRegion string
	var strLocationCountry string
	var strLocationMonth string
	var slotData map[string]alexa.Slot
	var strLocationCitySlot alexa.Slot
	var strLocationStateSlot alexa.Slot
	var strLocationRegionSlot alexa.Slot
	var strLocationCountrySlot alexa.Slot
	var strLocationMonthSlot alexa.Slot

	if resumingPrior {

		eventData = sessionData.Eventdata
		strLocationCity = sessionData.Name

	} else {
		slotData = request.Body.Intent.Slots

		strLocationCitySlot = slotData[LOCATION_CITY_SLOT]
		strLocationCity = strings.Title(strLocationCitySlot.Value)
		strLocationStateSlot = slotData[LOCATION_STATE_SLOT]
		strLocationState = strings.Title(strLocationStateSlot.Value)
		strLocationRegionSlot = slotData[LOCATION_REGION_SLOT]
		strLocationRegion = strings.Title(strLocationRegionSlot.Value)
		strLocationCountrySlot = slotData[LOCATION_COUNTRY_SLOT]
		strLocationCountry = strings.Title(strLocationCountrySlot.Value)
		strLocationMonthSlot = slotData[LOCATION_MONTH_SLOT]
		strLocationMonth = strings.Title(strLocationMonthSlot.Value)

		var err error
		strLocationCity = strings.Replace(strLocationCity, ",", "", -1) //remove random occurrence of a captured , char

		if strLocationRegion != "" {
			//This attempts to resolve the issue for English (IN), which is the only locale that does not provide for AMAZON.US_STATE
			//Region in many cases includes states, so if the language is English (IN), and Region slot data is present, use it as data for State
			strLocationState = strLocationRegion
		}
		fmt.Println("incoming slot info: City: " + strLocationCity + ", State: " + strLocationState + ", Country: " + strLocationCountry + ", Month: " + strLocationMonth)


		//---- Perform the Fetch of Event Data for the Artist
		eventData, err = fetchLocationData(strLocationCity, strLocationState, strLocationCountry, strLocationMonth)
		if err != nil {
			fmt.Println("Error received from fetchLocationData: ", err.Error())
		}

		var speechText string
		if eventData == nil || len(eventData) == 0 {
			speechText = "I couldn't find any events in " + strLocationCity
		} else {
			speechText = "Here are events happening near " + strLocationCity
		}

		if strLocationState != "" {
			speechText += ", " + strLocationState
		}
		if strLocationCountry != "" {
			speechText += ", " + strLocationCountry
		}
		if strLocationMonth != "" {
			speechText += " in " + strLocationMonth
		}

		primarySSMLText.Say(speechText)
		primarySSMLText.Pause("1000")
		cardTextContent += speechText + "\n\n"

	}

	//This variable is setup to hold APL custom Display property content
	customDisplayData := alexa.CustomDataToDisplay{
		ItemsListContent: make([]string, 3),
	}
	sessAttrData := make(map[string]interface{})
	titleString := ""

	if len(eventData) > 3 {
		for j := 0; j < 3; j++ {
			thisItem := eventData[j]

			primarySSMLText.Say(thisItem)
			primarySSMLText.Pause("1000")

			//This variable will store and be used to pass the text/content that needs to be displayed on the APL template
			customDisplayData.ItemsListContent[j] = thisItem
			cardTextContent += thisItem + "\n"
		}

		repromptString := "Would you like to hear more events?"
		primarySSMLText.Say(repromptString)
		primarySSMLText.Pause("1000")

		repromptSSMLText.Say(repromptString)
		repromptSSMLText.Pause("1000")

		//Save session attributes data for reentry, should the user answer yes to 'more' details..
		eventData = eventData[3:]

		sessionData.Eventdata = eventData
		sessionData.Intent = LOCATION_INTENT
		strCityStateCountry := strLocationCity
		if strLocationState != "" {
			strCityStateCountry += ", " + strLocationState
		}
		if strLocationCountry != "" {
			strCityStateCountry += ", " + strLocationCountry
		}

		sessionData.Name = strCityStateCountry
		sessionData.ID = LOCATION_ID
		sessAttrData["dataToSave"] = sessionData

		titleString = "Upcoming events near " + strCityStateCountry

	} else {

		//Is there at least 1 event left?
		if len(eventData) > 0 {

			for idx, item := range eventData {
				primarySSMLText.Say(item)
				primarySSMLText.Pause("1000")
				customDisplayData.ItemsListContent[idx] = item
				cardTextContent += item + "\n"
			}

			var textToSay = "Please ask another question like, Who is coming to Staples Center,  Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida. Say Cancel to exit. "
			primarySSMLText.Say("There are no additional events.")
			primarySSMLText.Pause("1000")
			primarySSMLText.Say(textToSay)
			cardTextContent += "There are no additional events.\n"
			repromptSSMLText.Say(textToSay)

			titleString = "Upcoming events near " + strLocationCity
			if strLocationState != "" {
				titleString += ", " + strLocationState
			}
			if strLocationCountry != "" {
				titleString += ", " + strLocationCountry
			}

		} else {
			//Couldn't find at least one event.. so either the Value provided was bad, OR the value was in fact good, but there are no events.
			//In either case, shoot off an SNS for research..
			err := clients.PublishSNSMessage(os.Getenv("SNS_TOPIC"), "Music Man Notification", "Music Man user request failure for LocationIntent, City: "+strLocationCity+", State: "+strLocationState+", Country: "+strLocationCountry+", Month: "+strLocationMonth)
			if err != nil {
				fmt.Println("Error sending SNS notification message")
			}

			//If both state and country were not included, then likely there are no events because we couldn't figure out a location id. Ask the
			//user to try again by providing an example.
			if strLocationState == "" && strLocationCountry == "" {
				primarySSMLText.Say("This may be due to not also providing either a U.S. state or a country along with your City.")
				primarySSMLText.Pause("1000")
				primarySSMLText.Say("Please re-ask your question similar to this, What is happening in Las Vegas, Nevada.")
				primarySSMLText.Pause("500")
				repromptSSMLText.Say("Please re-ask your question similar to this, What is happening in Las Vegas, Nevada.")
				repromptSSMLText.Pause("1000")
			} else {
				primarySSMLText.Say("If you would like to ask another question, try one of these:")
				primarySSMLText.Pause("500")
				primarySSMLText.Say("Who is playing at Staples Center, Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida. You can say Cancel to exit. ")
				primarySSMLText.Pause("1000")
				repromptSSMLText.Say("Please ask another question like, Who is playing at Staples Center, Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida. Say Cancel to exit. ")
				repromptSSMLText.Pause("1000")
			}

			titleString = "There are no upcoming events near " + strLocationCity
			if strLocationState != "" {
				titleString += " " + strLocationState
			}
			if strLocationCountry != "" {
				titleString += " " + strLocationCountry
			}

		}

	}

	if alexa.SupportsAPL(&request) {

		customDisplayData.ArtistVenueImgURL = "NA"

		return alexa.NewAPLAskResponse(titleString,
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardTextContent,
			false,
			&sessAttrData,
			"Main",
			&customDisplayData)
	} else {
		return alexa.NewSimpleAskResponse(titleString,
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardTextContent,
			false,
			&sessAttrData)
	}

}
