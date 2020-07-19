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
	APIRequestArtistID            func(string) (*models.ArtistIDResponse, error)
	APIRequestArtistEventCalendar func(string) (*models.CalendarResponse, error)
	ARTIST_ID                     string
)

func init() {
	APIRequestArtistID = clients.APIRequestArtistID
	APIRequestArtistEventCalendar = clients.APIRequestEventCalendar
	ARTIST_ID = "NA" //This is a default value to use in the APL template. It will indicate to the APL to just use the default image as there isnt an event/artist image available
}

const (
	SongkickArtistImageURL = "https://images.sk-static.com/images/media/profile_images/artists/ARTISTID/huge_avatar"
	ARTIST_NAME_SLOT       = "artist"
	ARTIST_MONTH_SLOT      = "month"
	ARTIST_INTENT          = "ArtistIntent"
)

type ArtistData struct {
	ID        string
	Name      string
	Eventdata []string
}

func fetchArtistData(artist, month string) ([]string, error) {

	thisMonth := strings.Title(month)

	urlToFetch, err := clients.ConstructURLRequest("ArtistQuery", artist, "", "")
	if err != nil {
		return nil, err
	}
	//Make an API call to Songkick to get the ArtistID for this artist
	artistIDResponse, err := APIRequestArtistID(urlToFetch)
	if err != nil {
		return nil, err
	}

	if artistIDResponse.ResultsPage.TotalEntries == 0 {
		//This artist wasnt found, so return immediately..
		return nil, nil
	}

	//With the ArtistID, construct the Songkick API Calendar request url
	ARTIST_ID = strconv.Itoa(artistIDResponse.ResultsPage.Results.Artist[0].ID)
	urlToFetch, err = clients.ConstructURLRequest("ArtistCalendar", ARTIST_ID, "", "")
	if err != nil {
		return nil, err
	}

	fmt.Println("URL being fetched: ", urlToFetch)
	//Make an API call to Songkick to get the Artist's Event Calendar
	artistCalendarResponse, err := APIRequestArtistEventCalendar(urlToFetch)
	if err != nil {
		return nil, err
	}

	counter := 0
	var itemsToSave []string

	for _, item := range artistCalendarResponse.ResultsPage.Results.Event {
		//If the user passed a Month as part of their request.. then filter out just those events..
		//The end result may be that no events are included.
		dateString := " on " + common.ConvertDate(item.Start.Date)
		thisLocation := item.Location.City
		if strings.Contains(thisLocation, ", US") {
			//Songkick uses the State abbreviation so convert it. The state is now the LAST two chars in this string..
			thisLocation = common.ConvertStateAbbreviation(thisLocation)
		}

		if thisMonth != "" && strings.Contains(dateString, " "+thisMonth+" ") || thisMonth == "" {
			displayEventString := " at " + item.Venue.DisplayName + dateString + " in " + thisLocation
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
func HandleArtistIntent(request alexa.Request, resumingPrior bool, sessionData models.SessionData) alexa.Response {

	var eventData []string
	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var cardTextContent string
	var strArtist string
	var strArtistMonth string
	var slotData map[string]alexa.Slot
	var strArtistSlot alexa.Slot
	var strArtistMonthSlot alexa.Slot

	if resumingPrior {

		eventData = sessionData.Eventdata
		strArtist = sessionData.Name

	} else {
		slotData = request.Body.Intent.Slots

		strArtistSlot = slotData[ARTIST_NAME_SLOT]
		strArtist = strArtistSlot.Value

		if len(slotData) > 1 {
			strArtistMonthSlot = slotData[ARTIST_MONTH_SLOT]
			strArtistMonth = strArtistMonthSlot.Value
		}

		var err error
		//---- See if there's a corrected value item (in the DynamoDB table) that we should use for the artist
		strArtist = common.CheckDynamoForCorrectedValue(strArtist)

		//---- Perform the Fetch of Event Data for the Artist
		eventData, err = fetchArtistData(strArtist, strArtistMonth)
		if err != nil {
			fmt.Println("Error received from fetchArtistData: ", err.Error())
		}

		var speechText string
		if eventData == nil || len(eventData) == 0 {
			speechText = "I couldn't find any events for " + strings.Title(strArtist)
			if strArtistMonth != "" {
				speechText += " in " + strArtistMonth
			}

		} else {
			speechText = "Here is where " + strings.Title(strArtist) + " is playing"
			if strArtistMonth != "" {
				speechText += " in " + strArtistMonth
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
		sessionData.Intent = ARTIST_INTENT
		sessionData.Name = strArtist
		sessionData.ID = ARTIST_ID
		sessAttrData["dataToSave"] = sessionData

		titleString = "Upcoming events for " + strings.Title(strArtist)

	} else {

		//Is there at least 1 event left?
		if len(eventData) > 0 {

			for idx, item := range eventData {
				primarySSMLText.Say(item)
				primarySSMLText.Pause("1000")
				customDisplayData.ItemsListContent[idx] = item
				cardTextContent += item + "\n"
			}
			primarySSMLText.Say("There are no additional events. Please ask another question like, Who is coming to Staples Center,  Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida. Say Cancel to exit. ")
			primarySSMLText.Pause("1000")
			cardTextContent += "There are no additional events.\n"

			titleString = "Upcoming events for " + strings.Title(strArtist)

		} else {
			//Couldn't find at least one event.. so either the Value provided was bad, OR the value was in fact good, but there are no events.
			//In either case, shoot off an SNS for research..
			err := clients.PublishSNSMessage(os.Getenv("SNS_TOPIC"), "Music Man Notification", "Music Man user request failure for ArtistIntent, Artist: "+strArtist+", Month: "+strArtistMonth)
			if err != nil {
				fmt.Println("Error sending SNS notification message")
			}
			primarySSMLText.Say("If you would like to ask another question, try one of these:")
			primarySSMLText.Pause("500")
			primarySSMLText.Say("Who is coming to Staples Center, Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida. You can say Cancel to exit. ")
			primarySSMLText.Pause("1000")

			titleString = "There are no upcoming events for " + strings.Title(strArtist)

		}
	}

	if alexa.SupportsAPL(&request) {

		customDisplayData.ArtistVenueImgURL = strings.Replace(SongkickArtistImageURL, "ARTISTID", ARTIST_ID, 1)

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
