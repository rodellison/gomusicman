package alexa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckSupportsAPLWhenInterfacePresent(t *testing.T) {

	theRequest := &Request{
		Version: "1.0",
		Session: Session{},
		Body: ReqBody{
			Type: "LaunchRequest",
		},
		Context: Context{
			System: System{
				Device: Device{
					SupportedInterfaces: SupportedInterfaces{
						APL: APL{
							Runtime: Runtime{MaxVersion: "1.3"},
						},
					},
				},
			},
		},
	}

	response := SupportsAPL(theRequest)
	assert.True(t, response, "APL should be present")

}

func TestCheckSupportsAPLWhenInterfaceNotPresent(t *testing.T) {

	theRequest := &Request{
		Version: "1.0",
		Session: Session{},
		Body: ReqBody{
			Type: "LaunchRequest",
		},
		Context: Context{
			System: System{
				Device: Device{
					SupportedInterfaces: SupportedInterfaces{},
				},
			},
		},
	}

	response := SupportsAPL(theRequest)
	assert.False(t, response, "APL should not be present")

}

func TestCheckIsEnglishWhenEnglishLocale(t *testing.T) {

	//theRequest := &Request{
	//	Version: "1.0",
	//	Session: Session{},
	//	Body:    ReqBody{
	//		Type: "LaunchRequest",
	//		Locale: "en-US",  //English US
	//	},
	//}

	var response bool

	response = IsEnglish(LocaleAmericanEnglish)
	assert.True(t, response)
	response = IsEnglish(LocaleIndianEnglish)
	assert.True(t, response)
	response = IsEnglish(LocaleBritishEnglish)
	assert.True(t, response)
	response = IsEnglish(LocaleCanadianEnglish)
	assert.True(t, response)
	response = IsEnglish(LocaleAustralianEnglish)
	assert.True(t, response)

}

func TestCheckIsEnglishWhenNOTEnglishLocale(t *testing.T) {

	var response bool

	response = IsEnglish(LocaleItalian)
	assert.False(t, response)
	response = IsEnglish(LocaleGerman)
	assert.False(t, response)
	response = IsEnglish(LocaleSpanishUS)
	assert.False(t, response)
	response = IsEnglish(LocaleJapanese)
	assert.False(t, response)
}
