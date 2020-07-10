package main

import (
	"encoding/json"
	"encoding/xml"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodellison/alexa-slick-dealer/alexa"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

//FeedResponse struct will be used as the type for unmarshalling external HTML data into
type FeedResponse struct {
	Channel struct {
		Item []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

//RequestFeed handles fetching external HTML site data and Unmarhalling to a struct that can be used later
//within the respective handler functions
func RequestFeed(mode string) (FeedResponse, error) {
	endpoint, _ := url.Parse("https://slickdeals.net/newsearch.php")
	queryParams := endpoint.Query()
	queryParams.Set("mode", mode)
	queryParams.Set("searcharea", "deals")
	queryParams.Set("searchin", "first")
	queryParams.Set("rss", "1")
	endpoint.RawQuery = queryParams.Encode()
	response, err := http.Get(endpoint.String())
	if err != nil {
		return FeedResponse{}, err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var feedResponse FeedResponse
		xml.Unmarshal(data, &feedResponse)
		return feedResponse, nil
	}
}

//Centralized function to steer incoming alexa requests to the appropirate handler function
func IntentDispatcher(request alexa.Request) alexa.Response {
	var response alexa.Response

	if request.Body.Type == "LaunchRequest" {
		return HandleLaunchIntent(request)
	}

	switch request.Body.Intent.Name {
	case "FrontpageDealIntent":
		response = HandleDealIntent(request, "Frontpage", false)
	case "PopularDealIntent":
		response = HandleDealIntent(request, "Popular", false)
	case alexa.YesIntent:
		response = HandleDealIntent(request, "", true)
	case alexa.NoIntent:
		response = HandleHelpIntent(request)
	case alexa.HelpIntent:
		response = HandleHelpIntent(request)
	case alexa.StopIntent:
		response = HandleStopIntent(request)
	case alexa.CancelIntent:
		response = HandleStopIntent(request)
	case alexa.FallbackIntent:
		response = HandleHelpIntent(request)
	}
	return response
}

//Handler function for the initial Launch Request - when someone says.. Alexa, open slick dealer.."
func HandleLaunchIntent(request alexa.Request) alexa.Response {

	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var response alexa.Response

	primarySSMLText.Say("Hello!, and welcome to Slick Dealer. You can ask a question like, What are the frontpage deals, or What are the popular deals.")
	primarySSMLText.Pause("1000")
	primarySSMLText.Say("For instructions on what you can say, please say help me.")
	repromptSSMLText.Say("Please ask a question like, What are the frontpage deals, or What are the popular deals.")

	if alexa.SupportsAPL(&request) {
		customDisplayData := alexa.CustomDataToDisplay{
			ItemsListContent: make([]string, 3),
		}
		customDisplayData.ItemsListContent[0] = "You can ask a question like, What are the front page deals, or What are the popular deals."

		response = alexa.NewAPLAskResponse("Welcome to Slick Dealer",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			"You can ask a question like, What are the front page deals, or What are the popular deals.",
			false,
			nil,
			"Home",
			customDisplayData)
	} else {
		response = alexa.NewSimpleAskResponse("Welcome to Slick Dealer",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			"You can ask a question like, What are the front page deals, or What are the popular deals.",
			false,
			nil)
	}

	return response

}

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
			feedResponse, _ = RequestFeed("frontpage")
			primarySSMLText.Say("Here are the current Front page deals:")
			cardTextContent += "Here are the current Front page deals:\n"
		} else {
			feedResponse, _ = RequestFeed("popdeals")
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

//Handler function for the Help Intent
func HandleHelpIntent(request alexa.Request) alexa.Response {
	var primarySSMLText alexa.SSMLBuilder
	var repromptSSMLText alexa.SSMLBuilder
	var response alexa.Response

	primarySSMLText.Say("OK, Here are some of the things you can ask:")
	primarySSMLText.Pause("1000")
	primarySSMLText.Say("What are the frontpage deals or,")
	primarySSMLText.Pause("500")
	primarySSMLText.Say("What are the popular deals.")

	repromptSSMLText.Say("Please ask a question like, What are the frontpage deals, or " +
		"What are the popular deals. Say Cancel if you'd like to quit.")
	repromptSSMLText.Pause("500")

	if alexa.SupportsAPL(&request) {

		//This variable is setup to hold APL Display content
		customDisplayData := alexa.CustomDataToDisplay{
			ItemsListContent: make([]string, 3),
		}
		customDisplayData.ItemsListContent[0] = "Please ask a question like, What are the Frontpage deals, or What are the Popular deals. You can also say 'Cancel' to exit."

		response = alexa.NewAPLAskResponse("Slick Dealer Help",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			"Please ask a question like, What are the Frontpage deals, or What are the Popular deals. You can also say 'Cancel' to exit.",
			false,
			nil,
			"Help",
			customDisplayData)
	} else {
		response = alexa.NewSimpleAskResponse("Slick Dealer Help",
			primarySSMLText.Build(),
			repromptSSMLText.Build(),
			"Please ask a question like, What are the Frontpage deals, or What are the Popular deals. You can also say 'Cancel' to exit.",
			false,
			nil)
	}

	return response
}

//Handler function for the Stop and Cancel Intent
func HandleStopIntent(request alexa.Request) alexa.Response {

	var primarySSMLText alexa.SSMLBuilder
	var response alexa.Response

	primarySSMLText.Say("Thanks and have a great day!, Goodbye.")
	response = alexa.NewSimpleTellResponse(os.Getenv("SkillTitle"),
		primarySSMLText.Build(),
		"Thanks and have a great day!, Goodbye.",
		true,
		nil)

	return response
}

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
