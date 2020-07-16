package handlers

import (
	"github.com/rodellison/gomusicman/alexa"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHandleHelpIntent(t *testing.T) {
	theRequest := &alexa.Request{
		Version: "1.0",
		Session: alexa.Session{
			Application: alexa.Application{
				ApplicationID: os.Getenv("AppARN"),
			},
		},
		Body: alexa.ReqBody{
			Type: "IntentRequest",
			Intent: alexa.Intent{
				Name:  "AMAZON.HelpIntent",
				Slots: nil,
			},
		},
		Context: alexa.Context{},
	}
	response := HandleHelpIntent(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}


func TestHandleHelpIntentWithAPL(t *testing.T) {

	context := alexa.Context{
		System: alexa.System{
			Application: alexa.Application{},
			User:        alexa.User{},
			Device: alexa.Device{
				DeviceID: "JustATest",
				SupportedInterfaces: alexa.SupportedInterfaces{
					APL: alexa.APL{
						alexa.Runtime{
							MaxVersion: "1.3",
						},
					},
				},
			},
			APIEndPoint:    "",
			APIAccessToken: "",
		},
	}

	theRequest := &alexa.Request{
		Version: "1.0",
		Session: alexa.Session{
			Application: alexa.Application{
				ApplicationID: os.Getenv("AppARN"),
			},
		},
		Body: alexa.ReqBody{
			Type: "IntentRequest",
			Intent: alexa.Intent{
				Name:  "AMAZON.HelpIntent",
				Slots: nil,
			},
		},
		Context: context,
	}

	alexa.FileToRead = "../apl_template_export.json"

	response := HandleHelpIntent(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}