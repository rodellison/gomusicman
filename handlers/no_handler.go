package handlers

import (
	"github.com/rodellison/gomusicman/alexa"
	"os"
)

//Handler function for the Help Intent
func HandleNoIntent(request alexa.Request) alexa.Response {

	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var response alexa.Response

	primarySSMLText.Say("OK, If you would like to ask another question, go ahead and ask now.")
	primarySSMLText.Pause("500")
	primarySSMLText.Say("Otherwise, just say Cancel to exit.")
	primarySSMLText.Pause("100")


	repromptSSMLText.Say("If you would like to ask another question, go ahead and ask now.")
	repromptSSMLText.Pause("500")
	repromptSSMLText.Say("Otherwise, just say Cancel to exit.")

	cardText := "Ask a question similar to one of these:\n\nWho is coming to the Mohawk?\nWhere is Iron Maiden playing in June?\nWHAT is happening in Fort Lauderdale, Florida?"
	sessAttrData := make(map[string]interface{})

	if alexa.SupportsAPL(&request) {

		//This variable is setup to hold APL Display content
		customDisplayData := alexa.CustomDataToDisplay{
			ItemsListContent: make([]string, 3),
		}
		customDisplayData.ItemsListContent[0] = "Alexa, Ask The Music Man:"
		customDisplayData.ItemsListContent[1] = "WHO is coming to {venue}<br/>WHERE is {artist} playing<br/>WHAT is happening in {city} + {state} or {country}<br/>"
		customDisplayData.ItemsListContent[2] = "You can also add 'in {month}' to any of the requests above for specific dates"
		customDisplayData.ArtistVenueImgURL = "NA"

		response = alexa.NewAPLAskResponse(os.Getenv("SkillTitle"),
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			&sessAttrData,
			"Main",
			&customDisplayData)
	} else {
		response = alexa.NewSimpleAskResponse(os.Getenv("SkillTitle"),
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			&sessAttrData)
	}

	return response
}
