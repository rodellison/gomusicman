package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodellison/gomusicman/alexa"
	"github.com/rodellison/gomusicman/handlers"
	"os"
)

var (
	//These var definitions to help with Mock testing
	StopCancelHandler, HelpHandler, LaunchHandler func (request alexa.Request) alexa.Response
)

func init() {
	//Assign the real handler functions here, but override when testing..
	StopCancelHandler = handlers.HandleStopIntent
	HelpHandler = handlers.HandleHelpIntent
	LaunchHandler = handlers.HandleLaunchIntent
}


//Centralized function to steer incoming alexa requests to the appropriate handler function
func IntentDispatcher(request alexa.Request) alexa.Response {
	var response alexa.Response

	if request.Body.Type == "LaunchRequest" {
		return LaunchHandler(request)
	}

	switch request.Body.Intent.Name {
	//case "FrontpageDealIntent":
	//	response = HandleDealIntent(request, "Frontpage", false)
	//case "PopularDealIntent":
	//	response = HandleDealIntent(request, "Popular", false)
	//case alexa.YesIntent:
	//	response = HandleDealIntent(request, "", true)
	//case alexa.NoIntent:
	//	response = HandleHelpIntent(request)
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


/*

//Handler function for the both the FrontPage and Popular Deal Intents.
//Parameters passed allow this function to accommodate both the initial request, as well as subsequent
//requests as a result of the user saying 'yes' for more data
func HandleDealIntent(request alexa.Request, dealType string, resumingPrior bool) alexa.Response {

	var feedResponse FeedResponse
	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var sessAttrData map[string]interface{}
	var cardTextContent string

	if resumingPrior {

		incomingSessionAttrs := request.Session.Attributes
		incomingData, _ := json.Marshal(incomingSessionAttrs["dataToSave"])
		json.Unmarshal(incomingData, &feedResponse)

	} else {

		if dealType == "Frontpage" {
			_ = APIRequest("frontpage", "")
			primarySSMLText.Say("Here are the current Front page deals:")
			cardTextContent += "Here are the current Front page deals:\n"
		} else {
			_ = APIRequest("popdeals", "")
			primarySSMLText.Say("Here are the current popular deals:")
			cardTextContent += "Here are the current popular deals:\n"
		}
		primarySSMLText.Pause("1000")
	}

	//This variable is setup to hold APL custom Display property content
	customDisplayData := alexa.CustomDataToDisplay{
		ItemsListContent: make([]string, 3),
	}

	if len(feedResponse.Channel.Item) > 3 {
		for j := 0; j < 3; j++ {
			thisItem := feedResponse.Channel.Item[j]
			primarySSMLText.Say(thisItem.Title)
			primarySSMLText.Pause("1000")

			//This variable will store and be used to pass the text/content that needs to be displayed on the APL template
			customDisplayData.ItemsListContent[j] = thisItem.Title
			cardTextContent += thisItem.Title + "\n"
		}

		repromptString := "Would you like to hear more deals"
		primarySSMLText.Say(repromptString)
		primarySSMLText.Pause("1000")

		repromptSSMLText.Say(repromptString)
		repromptSSMLText.Pause("1000")

		//Save session attributes data for reentry, should the user answer yes to 'more' details..
		feedResponse.Channel.Item = feedResponse.Channel.Item[3:]

		sessAttrData = make(map[string]interface{})
		sessAttrData["dataToSave"] = feedResponse

	} else {
		for idx, item := range feedResponse.Channel.Item {
			primarySSMLText.Say(item.Title)
			primarySSMLText.Pause("1000")
			customDisplayData.ItemsListContent[idx] = item.Title
			cardTextContent += item.Title + "\n"
		}
		sessAttrData = nil
		primarySSMLText.Say("There are no additional deals. Please ask another question like, What are the popular deals or What are the frontpage deals. Say Cancel to exit. ")
		primarySSMLText.Pause("1000")
	}

	if alexa.SupportsAPL(&request) {

		return alexa.NewAPLAskResponse(dealType+" Deals",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardTextContent,
			false,
			sessAttrData,
			"ItemsList",
			customDisplayData)
	} else {
		return alexa.NewSimpleAskResponse(dealType+" Deals",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			cardTextContent,
			false,
			sessAttrData)
	}

}

 */

//handler() is the first call from the lambda handler, first checking if the caller is a defined 'Alexa Skill - has an ARN we approve', and
//if so, will proceed
func Handler(request alexa.Request) (alexa.Response, error) {

	//Ensure this lambda function/code is invoked through the associated Alexa Skill, and not called directly
	if request.Session.Application.ApplicationID != os.Getenv("AppARN") {
		var primarybuilder alexa.SSMLBuilder
		primarybuilder.Say("Sorry, not authorized. Please enable and use this skill through an approved Alexa device.")
		return alexa.NewSimpleTellResponse("Not authorized",
			primarybuilder.Build(),
			"Not authorized, Please enable and use this skill through an approved Alexa device.",
			true,
			nil), nil
	} else {
		return IntentDispatcher(request), nil
	}

}

//main() is the entry point for the lambda handler
func main() {
	lambda.Start(Handler)
}
