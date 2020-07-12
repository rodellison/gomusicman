package main

import (
	"github.com/rodellison/gomusicman/alexa"
	"github.com/rodellison/gomusicman/clients"
	"github.com/rodellison/gomusicman/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)


func init() {
	clients.TheHTTPClient = &mocks.MockHTTPClient{}
}

func TestHandleLaunchIntent(t *testing.T) {

	theRequest := &alexa.Request{
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

	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty")
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

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
	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

func TestHandleFallback(t *testing.T) {

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
				Name:  "AMAZON.FallbackIntent",
				Slots: nil,
			},
		}, Context: alexa.Context{},
	}
	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

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
	response, _ := Handler(*theRequest)
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
	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

func TestNotAuthorized(t *testing.T) {

	theRequest := &alexa.Request{
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
	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.Contains(t, response.Body.OutputSpeech.SSML, "please enable and use this skill through an approved alexa device.")
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

	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

func TestHandleLaunchIntentWithAPL(t *testing.T) {

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
			Type: "LaunchRequest",
		},
		Context: context,
	}

	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}


/*
func TestHandleFrontpageDealIntent(t *testing.T) {

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
				Name:  "FrontpageDealIntent",
				Slots: nil,
			},
		},
		Context: alexa.Context{},
	}

	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty")
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

func TestHandlePopularDealIntent(t *testing.T) {
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
				Name:  "PopularDealIntent",
				Slots: nil,
			},
		},
		Context: alexa.Context{},
	}

	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty")
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}



func TestHandleDealResumeDetails(t *testing.T) {
	sessionAttrData := make(map[string]interface{})
	sessionAttrData["dataToSave"] = "some data"

	theRequest := &alexa.Request{
		Version: "1.0",
		Session: alexa.Session{
			Application: alexa.Application{
				ApplicationID: os.Getenv("AppARN"),
			},
			Attributes: sessionAttrData,
		},
		Body: alexa.ReqBody{
			Type: "IntentRequest",
			Intent: alexa.Intent{
				Name:  "AMAZON.YesIntent",
				Slots: nil,
			},
		},
		Context: alexa.Context{},
	}
	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}

func TestHandleNoIntent(t *testing.T) {

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
				Name:  "AMAZON.NoIntent",
				Slots: nil,
			},
		},
		Context: alexa.Context{},
	}
	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}


func TestHandlePopularDealsIntentWithAPL(t *testing.T) {

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
				Name:  "PopularDealIntent",
				Slots: nil,
			},
		},
		Context: context,
	}

	response, _ := Handler(*theRequest)
	assert.NotEmpty(t, response, "The response should not be empty", false)
	assert.NotEmpty(t, response.Body.OutputSpeech, "There should be output speech")
}


*/

