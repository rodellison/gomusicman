package alexa

import (
	"os"
	"strings"
)

// Modified version of Arien Malec's work
// https://github.com/arienmalec/alexa-go
// https://medium.com/@amalec/alexa-skills-with-go-54db0c21e758

var (
	//Options include Standard, Simple, LinkAccount, AskForPermissionsConsent
	CardTypeToUse = "Standard"
)

//The Response structure encapsulates session and body struct data that are returned to the user's device //to reply to their request.
type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              ResBody                `json:"response"`
}

//The Response Body structure is a core component of the Response, encapulating the the actual json components and data that hold the speech or display oriented data returned to the user's device to reply to their request.
type ResBody struct {
	OutputSpeech     *Payload     `json:"outputSpeech,omitempty"`
	Card             *Payload     `json:"card,omitempty"`
	Reprompt         *Reprompt    `json:"reprompt,omitempty"`
	Directives       []Directives `json:"directives,omitempty"`
	ShouldEndSession bool         `json:"shouldEndSession"`
}

//Reprompt is a small component struct for encapsulating the reprompt portion of speech that will be provided back to the user.
type Reprompt struct {
	OutputSpeech *Payload `json:"outputSpeech,omitempty"`
}

//Directives are a specialized struct for encapsulating data (audio, visual, etc.) that are provided back to the user.
//This isn't the actual output speech, but rather commands and necessary code for making content 'play' or 'display' on
//the end user device.
type Directives struct {
	Type          string         `json:"type,omitempty"`
	Token         string         `json:"token,omitempty"`
	SlotToElicit  string         `json:"slotToElicit,omitempty"`
	UpdatedIntent *UpdatedIntent `json:"UpdatedIntent,omitempty"`
	PlayBehavior  string         `json:"playBehavior,omitempty"`
	AudioItem     struct {
		Stream struct {
			Token                string `json:"token,omitempty"`
			URL                  string `json:"url,omitempty"`
			OffsetInMilliseconds int    `json:"offsetInMilliseconds,omitempty"`
		} `json:"stream,omitempty"`
	} `json:"audioItem,omitempty"`
	Document    APLDocument    `json:"document,omitempty"`
	DataSources APLDataSources `json:"datasources,omitempty"`
	TimeoutType string         `json:"timeoutType,omitempty"`
}

type UpdatedIntent struct {
	Name               string                 `json:"name,omitempty"`
	ConfirmationStatus string                 `json:"confirmationStatus,omitempty"`
	Slots              map[string]interface{} `json:"slots,omitempty"`
}

type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

//The Payload struct contains the key data that makes up Output speech, card content, images, etc. One to many of the
//properties may be included as needed.
type Payload struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	SSML    string `json:"ssml,omitempty"`
	Content string `json:"content,omitempty"`
	Image   Image  `json:"image,omitempty"`
}

//Response oriented functions ------------------------------------
func ParseString(text string) string {
	text = strings.ToLower(text)
	text = strings.Replace(text, "&", " and ", -1)
	text = strings.Replace(text, "+", " plus ", -1)

	text = strings.Replace(text, "AT&T", "a. t. and t", -1)
	text = strings.Replace(text, "BB&T", "b. b. and t", -1)
	text = strings.Replace(text, "US", "u. s.", -1)


	return text
}

//Helper function for constructing the Images component used for Response Cards
func getImages() Image {

	//Note: The actual image links (ENV, hardcoded, etc.) MUST be https, and not http
	images := &Image{
		SmallImageURL: os.Getenv("SmallImg"),
		LargeImageURL: os.Getenv("LargeImg"),
	}

	return *images
}

//NewSimpleTellResponse constructs a non reprompt oriented Response structure that can be returned to the Alexa user who is NOT using a
//display capable device. (i.e. Echo, Dot, Tap)
func NewSimpleTellResponse(title, ssmlPrimaryText, cardText string, endSession bool, sessionDataToSave *map[string]interface{}) Response {

	//This version is for non Display oriented Alexa devices (i.e.  Echo, Dot).
	r := Response{
		Version:           "1.0",
		SessionAttributes: *sessionDataToSave,
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "SSML",
				SSML: ssmlPrimaryText,
			},
			Card: &Payload{
				Type:    CardTypeToUse,
				Title:   title,
				Text:    cardText,
				Content: cardText,
				Image:   getImages(),
			},
			ShouldEndSession: endSession,
		},
	}

	return r
}

//NewSimpleAskResponse constructs a Reprompt oriented Response structure that can be returned to the Alexa user who is NOT using a
//display capable device. (i.e. Echo, Dot, Tap)
func NewSimpleAskResponse(title, ssmlPrimaryText, ssmlRepromptText, cardText string, endSession bool, sessionDataToSave *map[string]interface{}) Response {

	//This version is for non Display oriented Alexa devices (i.e.  Echo, Dot).
	r := Response{
		Version:           "1.0",
		SessionAttributes: *sessionDataToSave,
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "SSML",
				SSML: ssmlPrimaryText,
			},
			Reprompt: &Reprompt{
				OutputSpeech: &Payload{
					Type: "SSML",
					SSML: ssmlRepromptText,
				},
			},
			Card: &Payload{
				Type:    CardTypeToUse,
				Title:   title,
				Text:    cardText,
				Content: cardText,
				Image:   getImages(),
			},
			ShouldEndSession: endSession,
		},
	}

	return r
}

//NewAPLTellResponse constructs a non reprompt oriented Response structure that can be returned to the Alexa user who IS using a
//display capable device. (i.e. Show, Firestick)
func NewAPLTellResponse(title, ssmlPrimaryText, cardText string, endSession bool, sessionDataToSave *map[string]interface{}, layoutToUse string, contentToUse *CustomDataToDisplay) Response {

	//This version is for APL Display oriented Alexa devices (i.e.  Show, Firestick).
	myAPLDocData, err := FetchAPL()
	if err != nil {
		//If an APL (load or unmarshal error occurred), return the simple version response instead as a fallback
		return NewSimpleTellResponse(title, ssmlPrimaryText, cardText, endSession, sessionDataToSave)
	}

	//Now Adjust the Data source properties as needed
	//This sets which APL layout will be used for displaying content
	myAPLDocData.APLDataSources.TemplateData.Properties.LayoutToUse = layoutToUse

	switch layoutToUse {

	case "Home":
		myAPLDocData.APLDataSources.TemplateData.Properties.HeadingText = "Welcome to " + os.Getenv("SkillTitle")
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[0] = contentToUse.ItemsListContent[0]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventImageUrl = "NA"
		myAPLDocData.APLDataSources.TemplateData.Properties.HintString = "Where is Iron Maiden playing in July"

	case "Help":
		myAPLDocData.APLDataSources.TemplateData.Properties.HeadingText = title
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[0] = contentToUse.ItemsListContent[0]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[1] = contentToUse.ItemsListContent[1]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[2] = contentToUse.ItemsListContent[2]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventImageUrl = "NA"
		myAPLDocData.APLDataSources.TemplateData.Properties.HintString = "Who is coming to the Mohawk in May"

	case "Events":
		myAPLDocData.APLDataSources.TemplateData.Properties.HeadingText = title
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[0] = contentToUse.ItemsListContent[0]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[1] = contentToUse.ItemsListContent[1]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[2] = contentToUse.ItemsListContent[2]
	}

	APLDirective := make([]Directives, 1)
	APLDirective[0].Type = "Alexa.Presentation.APL.RenderDocument"
	APLDirective[0].TimeoutType = "SHORT"
	APLDirective[0].Document = myAPLDocData.APLDocument
	APLDirective[0].DataSources = myAPLDocData.APLDataSources

	r := Response{
		Version:           "1.0",
		SessionAttributes: *sessionDataToSave,
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "SSML",
				SSML: ssmlPrimaryText,
			},
			Directives: APLDirective,
			Card: &Payload{
				Type:    CardTypeToUse,
				Title:   title,
				Text:    cardText,
				Content: cardText,
				Image:   getImages(),
			},
			ShouldEndSession: endSession,
		},
	}

	return r
}

//NewAPLAskResponse constructs a Reprompt oriented Response structure that can be returned to the Alexa user who IS using a
//display capable device. (i.e. Show, Firestick)
func NewAPLAskResponse(title, ssmlPrimaryText, ssmlRepromptText, cardText string, endSession bool, sessionDataToSave *map[string]interface{}, layoutToUse string, contentToUse *CustomDataToDisplay) Response {

	//This version is for APL Display oriented Alexa devices (i.e.  Show, Firestick).
	myAPLDocData, err := FetchAPL()
	if err != nil {
		//If an APL (load or unmarshal error occurred), return the simple version response instead as a fallback
		return NewSimpleTellResponse(title, ssmlPrimaryText, cardText, endSession, sessionDataToSave)
	}

	//Now Adjust the Data source properties as needed
	//This sets which APL layout will be used for displaying content
	myAPLDocData.APLDataSources.TemplateData.Properties.LayoutToUse = layoutToUse

	switch layoutToUse {

	case "Home":
		myAPLDocData.APLDataSources.TemplateData.Properties.HeadingText = "Welcome to " + os.Getenv("SkillTitle")
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[0] = contentToUse.ItemsListContent[0]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventImageUrl = "NA"
		myAPLDocData.APLDataSources.TemplateData.Properties.HintString = "Where is Iron Maiden playing in July"
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[0] = contentToUse.ItemsListContent[0]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[1] = contentToUse.ItemsListContent[1]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[2] = contentToUse.ItemsListContent[2]
	case "Help":
		myAPLDocData.APLDataSources.TemplateData.Properties.HeadingText = title
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[0] = contentToUse.ItemsListContent[0]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[1] = contentToUse.ItemsListContent[1]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[2] = contentToUse.ItemsListContent[2]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventImageUrl = "NA"
		myAPLDocData.APLDataSources.TemplateData.Properties.HintString = "Who is coming to the Mohawk in May"

	case "Events":
		myAPLDocData.APLDataSources.TemplateData.Properties.HeadingText = title
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[0] = contentToUse.ItemsListContent[0]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[1] = contentToUse.ItemsListContent[1]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventText[2] = contentToUse.ItemsListContent[2]
		myAPLDocData.APLDataSources.TemplateData.Properties.EventImageUrl = contentToUse.ArtistImgURL
	}

	APLDirective := make([]Directives, 1)
	APLDirective[0].Type = "Alexa.Presentation.APL.RenderDocument"
	APLDirective[0].TimeoutType = "SHORT"
	APLDirective[0].Document = myAPLDocData.APLDocument
	APLDirective[0].DataSources = myAPLDocData.APLDataSources

	r := Response{
		Version:           "1.0",
		SessionAttributes: *sessionDataToSave,
		Body: ResBody{
			OutputSpeech: &Payload{
				Type: "SSML",
				SSML: ssmlPrimaryText,
			},
			Reprompt: &Reprompt{
				OutputSpeech: &Payload{
					Type: "SSML",
					SSML: ssmlRepromptText,
				},
			},
			Directives: APLDirective,
			Card: &Payload{
				Type:    CardTypeToUse,
				Title:   title,
				Text:    cardText,
				Content: cardText,
				Image:   getImages(),
			},
			ShouldEndSession: endSession,
		},
	}

	return r
}
