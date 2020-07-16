package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodellison/gomusicman/alexa"
	"github.com/rodellison/gomusicman/handlers"
	"github.com/rodellison/gomusicman/models"
	"os"
)

var (
	//These var definitions to help with Mock testing
	StopCancelHandler, HelpHandler, LaunchHandler func(alexa.Request) alexa.Response
	ArtistHandler, VenueHandler                   func(alexa.Request, bool, models.SessionData) alexa.Response
)

func init() {
	//Assign the real handler functions here, but override when testing..
	StopCancelHandler = handlers.HandleStopIntent
	HelpHandler = handlers.HandleHelpIntent
	LaunchHandler = handlers.HandleLaunchIntent
	ArtistHandler = handlers.HandleArtistIntent
	//	VenueHandler = handlers.HandleVenueIntent
}

//Centralized function to steer incoming alexa requests to the appropriate handler function
func IntentDispatcher(request alexa.Request) alexa.Response {

	var response alexa.Response
	var sessionData models.SessionData

	if request.Body.Type == "LaunchRequest" {
		return LaunchHandler(request)
	}

	switch request.Body.Intent.Name {
	case "ArtistIntent":
		response = ArtistHandler(request, false, sessionData)
	case "VenueIntent":
		response = VenueHandler(request, false, sessionData)
	case alexa.YesIntent:

		incomingSessionAttrs := request.Session.Attributes
		incomingData, _ := json.Marshal(incomingSessionAttrs["dataToSave"])
		json.Unmarshal(incomingData, &sessionData)
		if len(sessionData.Eventdata) == 0 {
			response = HelpHandler(request)
		} else {
			//because we've unmashalled the session data already, rather than do it again inside the respective handlers,
			//just pass it along as a parm.
			if sessionData.Intent == "ArtistIntent" {
				response = ArtistHandler(request, true, sessionData)
			} else {
				response = VenueHandler(request, true, sessionData)
			}
		}

	case alexa.NoIntent:
		response = HelpHandler(request)
	case alexa.StopIntent:
		response = StopCancelHandler(request)
	case alexa.CancelIntent:
		response = StopCancelHandler(request)
	case alexa.FallbackIntent:
		response = HelpHandler(request)
	case alexa.HelpIntent:
		response = HelpHandler(request)
	}
	return response
}

//handler() is the first call from the lambda handler, first checking if the caller is coming from an expected 'Alexa Skill ARN.
//if so, proceed, if not - send not auth response
func Handler(request alexa.Request) (alexa.Response, error) {

	//Ensure this lambda function/code is invoked through the associated Alexa Skill, and not called directly
	if request.Session.Application.ApplicationID != os.Getenv("AppARN") {
		var primarybuilder alexa.SSMLBuilder
		primarybuilder.Say("Sorry, not authorized. Please enable and use this skill through an approved Alexa device.")

		sessAttrData := make(map[string]interface{})
		return alexa.NewSimpleTellResponse("Not authorized",
			primarybuilder.Build(),
			"Not authorized, Please enable and use this skill through an approved Alexa device.",
			true,
			&sessAttrData), nil
	} else {
		return IntentDispatcher(request), nil
	}

}

//main() is the entry point for the lambda handler
func main() {
	lambda.Start(Handler)
}
