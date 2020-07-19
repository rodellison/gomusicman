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
	primarySSMLText.Say("Hello!, The Music Man can tell you who's coming to a particular venue, where an artist is playing, and what concerts are happening in a particular city.")
	primarySSMLText.Pause("300")
	primarySSMLText.Say("Try asking a question similar to one of these:")
	primarySSMLText.Pause("200")
	primarySSMLText.Say("Who is coming to Staples Center, or Where is Iron Maiden playing, or What is happening in Fort Lauderdale, Florida")

	repromptSSMLText.Say("Please ask a question similar to one of these")
	repromptSSMLText.Pause("200")
	repromptSSMLText.Say("Who is coming to Staples Center, or Where is Iron Maiden playing in July, or What is happening in Fort Lauderdale, Florida?")

	cardText := "Try asking a question similar to one of these: \n Who is coming to Staples Center?\n Where is Iron Maiden playing?\n What is happening in Fort Lauderdale, Florida?"
	sessAttrData := make(map[string]interface{})

	if alexa.SupportsAPL(&request) {
		customDisplayData := alexa.CustomDataToDisplay{
			ItemsListContent: make([]string, 3),
		}

		customDisplayData.ItemsListContent[0] = "WHO is coming to Staples Center?"
		customDisplayData.ItemsListContent[1] = "WHERE is Iron Maiden playing?"
		customDisplayData.ItemsListContent[2] = "WHAT is happening in Fort Lauderdale, Florida"

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
