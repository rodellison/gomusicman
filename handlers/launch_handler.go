package handlers

import (
	"github.com/rodellison/gomusicman/alexa"
	"os"
)

//Handler function for the initial Launch Request - when someone says.. Alexa, open slick dealer.."
func HandleLaunchIntent(request alexa.Request) alexa.Response {

	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var response alexa.Response

	primarySSMLText.Say("<audio src='soundbank://soundlibrary/musical/amzn_sfx_musical_drone_intro_02'/>")
	primarySSMLText.Say("Hello!, The Music Man can tell you where an artist is playing, or who's coming to a particular venue.")
	primarySSMLText.Pause("300")
	primarySSMLText.Say("Try asking a question similar to one of these:")
	primarySSMLText.Pause("200")
	primarySSMLText.Say("Who is coming to Staples Center, or Where is Iron Maiden playing in July?")

	repromptSSMLText.Say("Please ask a question similar to one of these")
	repromptSSMLText.Pause("200")
	repromptSSMLText.Say("Who is coming to Staples Center, or Where is Iron Maiden playing in July?")

	cardText := "Try asking a question similar to one of these: \n Who is coming to Staples Center, or Where is Iron Maiden playing in July?"
	sessAttrData := make(map[string]interface{})

	if alexa.SupportsAPL(&request) {
		customDisplayData := alexa.CustomDataToDisplay{
			ItemsListContent: make([]string, 3),
		}

		customDisplayData.ItemsListContent[0] = "Ask a question similar to one of these:"
		customDisplayData.ItemsListContent[1] = "Who is coming to Staples Center?"
		customDisplayData.ItemsListContent[2] = "Where is Iron Maiden playing"

		response = alexa.NewAPLAskResponse("Welcome to "+os.Getenv("SkillTitle"),
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			&sessAttrData,
			"Home",
			&customDisplayData)
	} else {
		response = alexa.NewSimpleAskResponse("Welcome to "+os.Getenv("SkillTitle"),
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardText,
			false,
			&sessAttrData)
	}

	return response

}
