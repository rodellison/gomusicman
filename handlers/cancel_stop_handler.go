package handlers

import (
	"github.com/rodellison/gomusicman/alexa"
	"os"
)

//Handler function for the Stop and Cancel Intent
func HandleStopIntent(request alexa.Request) alexa.Response {

	var primarySSMLText alexa.SSMLBuilder
	var response alexa.Response

	primarySSMLText.Say("Thanks for using the Music Man. Goodbye.")
	primarySSMLText.Say("<audio src='soundbank://soundlibrary/musical/amzn_sfx_musical_drone_intro_02'/>")
	response = alexa.NewSimpleTellResponse(os.Getenv("SkillTitle"),
		primarySSMLText.Build(),
		"Thanks for using the Music Man!",
		true,
		nil)

	return response
}

