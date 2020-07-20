package main

import (
	"fmt"
	"github.com/rodellison/gomusicman/alexa"
	"github.com/rodellison/gomusicman/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	theRequest  alexa.Request
	theResponse alexa.Response
)

func init() {
	//Set up some dummy request and response objects with minimal data
	theResponse = alexa.Response{}
	theRequest = alexa.Request{
		Version: "1.0",
		Session: alexa.Session{
			Application: alexa.Application{
				ApplicationID: os.Getenv("AppARN"),
			},
		},
		Body: alexa.ReqBody{
			Type: "LaunchRequest",
		},
		Context: alexa.Context{},
	}

	//define some dummy handlers for each Intent, and just return dummy response
	//Testing for the actual handlers will be performed in their own respective test files

	LaunchHandler = func(alexa.Request) alexa.Response {
		return theResponse
	}
	StopCancelHandler = func(alexa.Request) alexa.Response {
		return theResponse
	}
	HelpHandler = func(alexa.Request) alexa.Response {
		return theResponse
	}
	ArtistHandler = func(alexa.Request, bool, models.SessionData) alexa.Response {
		return theResponse
	}
	VenueHandler = func(alexa.Request, bool, models.SessionData) alexa.Response {
		return theResponse
	}
	NoHandler = func(alexa.Request) alexa.Response {
		return theResponse
	}

}

func TestNotAuthorized(t *testing.T) {

	theNotAuthRequest := &alexa.Request{
		Version: "1.0",
		Session: alexa.Session{
			Application: alexa.Application{
				ApplicationID: "SomeIncorrectOrNOARNData",
			},
		},
		Body: alexa.ReqBody{
			Type: "LaunchRequest",
		},
		Context: alexa.Context{},
	}
	response, _ := Handler(*theNotAuthRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.Contains(t, response.Body.OutputSpeech.SSML, "Please enable and use this skill through an approved Alexa device.")
}

func TestIntentDispatcher(t *testing.T) {

	fmt.Println("Test Launch Handler through IntentDispatcher")
	response, _ := Handler(theRequest)
	assert.Equal(t, theResponse, response)

	fmt.Println("Test Help Handler through IntentDispatcher")
	theRequest.Body.Type = "Intent"
	theRequest.Body.Intent = alexa.Intent{
		Name:  "AMAZON.HelpIntent",
		Slots: nil,
	}
	response, _ = Handler(theRequest)
	assert.Equal(t, theResponse, response)

	fmt.Println("Test Cancel Stop Handler through IntentDispatcher")
	theRequest.Body.Type = "Intent"
	theRequest.Body.Intent = alexa.Intent{
		Name:  "AMAZON.StopIntent",
		Slots: nil,
	}
	response, _ = Handler(theRequest)
	assert.Equal(t, theResponse, response)

	fmt.Println("Test Artist Handler through IntentDispatcher")
	theRequest.Body.Type = "Intent"
	theRequest.Body.Intent = alexa.Intent{
		Name:  "ArtistIntent",
		Slots: nil,
	}
	response, _ = Handler(theRequest)
	assert.Equal(t, theResponse, response)

	fmt.Println("Test FallbackHandler through IntentDispatcher")
	theRequest.Body.Type = "Intent"
	theRequest.Body.Intent = alexa.Intent{
		Name:  "AMAZON.FallbackIntent",
		Slots: nil,
	}
	response, _ = Handler(theRequest)
	assert.Equal(t, theResponse, response)

	fmt.Println("Test No Handler through IntentDispatcher")
	theRequest.Body.Type = "Intent"
	theRequest.Body.Intent = alexa.Intent{
		Name:  "AMAZON.NoIntent",
		Slots: nil,
	}
	response, _ = Handler(theRequest)
	assert.Equal(t, theResponse, response)

	//todo: test YES Handler

}
