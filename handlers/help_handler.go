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

	primarySSMLText.Say("OK, here are some ways you can ask for events from the Music Man.")
	primarySSMLText.Pause("500")
	primarySSMLText.Say("If you want to know upcoming events for a particular venue, say your question like this:")
	primarySSMLText.Pause("100")
	primarySSMLText.Say("Who is coming TO Staples Center, or Who is coming TO the Mohawk.")
	primarySSMLText.Pause("1000")
	primarySSMLText.Say("If you want to know where an artist is playing, try saying:")
	primarySSMLText.Pause("100")
	primarySSMLText.Say("Where is Iron Maiden playing, or ")
	primarySSMLText.Pause("500")
	primarySSMLText.Say("If you want to know events happening in a particular city, try saying:")
	primarySSMLText.Pause("100")
	primarySSMLText.Say("What is happening IN Fort Lauderdale, Florida.")
	primarySSMLText.Pause("1000")
	primarySSMLText.Say("You can narrow the events provided to a particular month as well. Try saying: ")
	primarySSMLText.Pause("100")
	primarySSMLText.Say("Where is Iron Maiden playing in July.")

	repromptSSMLText.Say("Ask a question similar to one of these: ")
	repromptSSMLText.Pause("500")
	repromptSSMLText.Say("Who is coming to the Mohawk, or Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida?")

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

		response = alexa.NewAPLAskResponse(os.Getenv("SkillTitle")+" Help",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			&sessAttrData,
			"Main",
			&customDisplayData)
	} else {
		response = alexa.NewSimpleAskResponse(os.Getenv("SkillTitle")+" Help",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			&sessAttrData)
	}

	return response
}
