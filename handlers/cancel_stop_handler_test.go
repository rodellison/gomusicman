package handlers

import (
	"github.com/rodellison/gomusicman/alexa"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHandleStopIntent(t *testing.T) {
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
				Name:  "AMAZON.StopIntent",
				Slots: nil,
			},
		}, Context: alexa.Context{},
	}
	response := HandleStopIntent(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

func TestHandleCancelIntent(t *testing.T) {

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
				Name:  "AMAZON.CancelIntent",
				Slots: nil,
			},
		}, Context: alexa.Context{},
	}
	response := HandleStopIntent(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}