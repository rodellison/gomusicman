package handlers

import (
	"github.com/rodellison/gomusicman/alexa"
	"os"
)

//Handler function for the Help Intent
func HandleHelpIntent(request alexa.Request) alexa.Response {
	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var response alexa.Response

	primarySSMLText.Say("There's various ways to ask for Artist or Venue information.")
	primarySSMLText.Pause("500")
	primarySSMLText.Say("Ask a question similar to one of these:")
	primarySSMLText.Pause("500")
	primarySSMLText.Say("Who is coming to the Mohawk, or ")
	primarySSMLText.Pause("100")
	primarySSMLText.Say("Where are the Rolling Stones playing in June. ")
	primarySSMLText.Pause("500")
	primarySSMLText.Say("To quit, just say cancel.")

	repromptSSMLText.Say("Ask a question similar to one of these: ")
	repromptSSMLText.Pause("500")
	repromptSSMLText.Say("Who is coming to the Mohawk, or Where are the Rolling Stones playing in June?")

	cardText := "Ask a question similar to one of these: \n Who is coming to the Mohawk, or Where are the Rolling Stones playing in June"

	if alexa.SupportsAPL(&request) {

		//This variable is setup to hold APL Display content
		customDisplayData := alexa.CustomDataToDisplay{
			ItemsListContent: make([]string, 3),
		}
		customDisplayData.ItemsListContent[0] = "Alexa, Ask The Music Man:"
		customDisplayData.ItemsListContent[1] = "When is {artist} playing<br/>" + "Who is coming to {venue}<br/>"
		customDisplayData.ItemsListContent[2] = "You can also add 'in {month}' to any of the requests above for specific dates"

		response = alexa.NewAPLAskResponse(os.Getenv("SkillTitle") + " Help",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			nil,
			"Help",
			customDisplayData)
	} else {
		response = alexa.NewSimpleAskResponse(os.Getenv("SkillTitle") + " Help",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			nil)
	}

	return response
}
