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
	APIRequestVenueID            func(string) (*models.VenueIDResponse, error)
	APIRequestVenueEventCalendar func(string) (*models.CalendarResponse, error)
	VENUE_ID                     string
)

func init() {
	APIRequestVenueID = clients.APIRequestVenueID
	APIRequestVenueEventCalendar = clients.APIRequestEventCalendar
	VENUE_ID = "NA" //This is a default value to use in the APL template. It will indicate to the APL to just use the default image as there isnt an event/Venue image available
}

const (
	SongkickVenueImageURL = "https://images.sk-static.com/images/media/profile_images/venues/VENUEID/huge_avatar"
	VENUE_NAME_SLOT       = "venue"
	VENUE_MONTH_SLOT      = "month"
	VENUE_INTENT          = "VenueIntent"
)

func fetchVenueData(venue, month string) ([]string, error) {

	thisMonth := strings.Title(month)

	urlToFetch, err := clients.ConstructURLRequest("VenueQuery", venue, "", "")
	if err != nil {
		return nil, err
	}
	//Make an API call to Songkick to get the VenueID for this venue
	venueIDResponse, err := APIRequestVenueID(urlToFetch)
	if err != nil {
		return nil, err
	}

	if venueIDResponse.ResultsPage.TotalEntries == 0 {
		//This venue wasnt found, so return immediately..
		return nil, nil
	}

	//With the VenueID, construct the Songkick API Calendar request url
	VENUE_ID = strconv.Itoa(venueIDResponse.ResultsPage.Results.Venue[0].ID)
	urlToFetch, err = clients.ConstructURLRequest("VenueCalendar", VENUE_ID, "", "")
	if err != nil {
		return nil, err
	}

	fmt.Println("URL being fetched: ", urlToFetch)
	//Make an API call to Songkick to get the Venue's Event Calendar
	venueCalendarResponse, err := APIRequestVenueEventCalendar(urlToFetch)
	if err != nil {
		return nil, err
	}

	counter := 0
	var itemsToSave []string

	for _, item := range venueCalendarResponse.ResultsPage.Results.Event {
		//If the user passed a Month as part of their request.. then filter out just those events..
		//The end result may be that no events are included.
		dateString := " on " + common.ConvertDate(item.Start.Date)
		if thisMonth != "" && strings.Contains(dateString, " "+thisMonth+" ") || thisMonth == "" {

			var displayEventString string
			if strings.Contains(item.DisplayName, " at ") {
				displayLocAt := strings.Index(item.DisplayName, " at ")
				displayEventString = item.DisplayName[0:displayLocAt] + dateString
			} else {
				displayEventString = item.DisplayName + dateString
			}

			if strings.Contains(item.DisplayName, "CANCELLED") {
				displayEventString += ", is CANCELLED."
			}
			itemsToSave = append(itemsToSave, displayEventString)
			counter += 1
		}
	}

	return itemsToSave, nil

}

//Parameters passed allow this function to accommodate both the initial request, as well as subsequent
//requests as a result of the user saying 'yes' for more data
func HandleVenueIntent(request alexa.Request, resumingPrior bool, sessionData models.SessionData) alexa.Response {

	var eventData []string
	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var cardTextContent string
	var strVenue string
	var strVenueMonth string
	var slotData map[string]alexa.Slot
	var strVenueSlot alexa.Slot
	var strVenueMonthSlot alexa.Slot

	if resumingPrior {

		eventData = sessionData.Eventdata
		strVenue = sessionData.Name

	} else {
		slotData = request.Body.Intent.Slots

		strVenueSlot = slotData[VENUE_NAME_SLOT]
		strVenue = strVenueSlot.Value

		if len(slotData) > 1 {
			strVenueMonthSlot = slotData[VENUE_MONTH_SLOT]
			strVenueMonth = strVenueMonthSlot.Value
		}

		var err error
		//---- See if there's a corrected value item (in the DynamoDB table) that we should use for the Venue

		strVenue = common.CheckDynamoForCorrectedValue(strVenue)
		fmt.Println("Returning from DynamoDB check, value to use: ", strVenue)

		//---- Perform the Fetch of Event Data for the Venue
		eventData, err = fetchVenueData(strVenue, strVenueMonth)
		if err != nil {
			fmt.Println("Error received from fetchVenueData: ", err.Error())
		}

		var speechText string
		if eventData == nil || len(eventData) == 0 {
			speechText = "I couldn't find any events at " + strings.Title(strVenue)
			if strVenueMonth != "" {
				speechText += " in " + strVenueMonth
			}

		} else {
			speechText = "Here are upcoming events at " + strings.Title(strVenue)
			if strVenueMonth != "" {
				speechText += " in " + strVenueMonth
			}
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
		sessionData.Intent = VENUE_INTENT
		sessionData.Name = strVenue
		sessionData.ID = VENUE_ID
		sessAttrData["dataToSave"] = sessionData

		titleString = "Upcoming events at " + strings.Title(strVenue)

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
			primarySSMLText.Pause("1000")
			cardTextContent += "There are no additional events.\n"

			repromptSSMLText.Say(textToSay)
			repromptSSMLText.Pause("1000")


			titleString = "Upcoming events at " + strings.Title(strVenue)

		} else {
			//Couldn't find at least one event.. so either the Value provided was bad, OR the value was in fact good, but there are no events.
			//In either case, shoot off an SNS for research..
			err := clients.PublishSNSMessage(os.Getenv("SNS_TOPIC"), "Music Man Notification", "Music Man user request failure for VenueIntent, Venue: "+strVenue+", Month: "+strVenueMonth)
			if err != nil {
				fmt.Println("Error sending SNS notification message")
			}
			primarySSMLText.Say("If you would like to ask another question, try one of these:")
			primarySSMLText.Pause("500")
			primarySSMLText.Say("Who is coming to Staples Center,  Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida. You can say Cancel to exit. ")

			repromptSSMLText.Say("Try asking a question like one of these:")
			repromptSSMLText.Pause("500")
			repromptSSMLText.Say("Who is coming to Staples Center,  Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida. You can say Cancel to exit. ")

			titleString = "There are no upcoming events at " + strings.Title(strVenue)

		}
	}

	if alexa.SupportsAPL(&request) {

		//customDisplayData.ArtistVenueImgURL = strings.Replace(SongkickVenueImageURL, "VENUEID", VENUE_ID, 1)
		//Hardcoding "NA" for now as the Venue images are very inconsistent from Songkick
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
