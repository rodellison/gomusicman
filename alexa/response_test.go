package alexa

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSimpleTellResponse(t *testing.T) {
	sessAttrData := make(map[string]interface{})
	response := NewSimpleTellResponse("TextTitle1", "<speak>TextToSpeak</speak>", "TextToDisplay", false, &sessAttrData)
	assert.Contains(t, response.Body.OutputSpeech.SSML, "<speak>TextToSpeak</speak>")
	assert.Equal(t, response.Body.Card.Text, "TextToDisplay")

}

func TestNewSimpleAskResponse(t *testing.T) {

	sessAttrData := make(map[string]interface{})
	response := NewSimpleAskResponse("TextTitle1", "<speak>TextToSpeak</speak>", "<speak>TextToReprompt</speak>", "TextToDisplay", false, &sessAttrData)
	assert.Contains(t, response.Body.OutputSpeech.SSML, "<speak>TextToSpeak</speak>")
	assert.Contains(t, response.Body.Reprompt.OutputSpeech.SSML, "<speak>TextToReprompt</speak>")
	assert.Equal(t, response.Body.Card.Text, "TextToDisplay")

}

func TestNewAPLTellResponse(t *testing.T) {

	FileToRead = "../apl_template_export.json"
	customDisplayData := CustomDataToDisplay{
		ItemsListContent: make([]string, 3),
	}
	customDisplayData.ItemsListContent[0] = "Test Content Filler"
	sessAttrData := make(map[string]interface{})

	response := NewAPLTellResponse("TextTitle1", "<speak>TextToSpeak</speak>", "TextToDisplay", false, &sessAttrData, "Home", &customDisplayData)
	assert.Equal(t, "Alexa.Presentation.APL.RenderDocument", response.Body.Directives[0].Type, "Alexa.Presentation.APL.RenderDocument")
	assert.Equal(t, "Home", response.Body.Directives[0].DataSources.TemplateData.Properties.LayoutToUse, "Home")
	response = NewAPLTellResponse("TextTitle1", "<speak>TextToSpeak</speak>", "TextToDisplay", false, &sessAttrData, "ItemsList", &customDisplayData)
	assert.Equal(t, "Alexa.Presentation.APL.RenderDocument", response.Body.Directives[0].Type, "Alexa.Presentation.APL.RenderDocument")
	assert.Equal(t, "ItemsList", response.Body.Directives[0].DataSources.TemplateData.Properties.LayoutToUse, "ItemsList")
	response = NewAPLTellResponse("TextTitle1", "<speak>TextToSpeak</speak>", "TextToDisplay", false, &sessAttrData, "Help", &customDisplayData)
	assert.Equal(t, "Alexa.Presentation.APL.RenderDocument", response.Body.Directives[0].Type, "Alexa.Presentation.APL.RenderDocument")
	assert.Equal(t, "Help", response.Body.Directives[0].DataSources.TemplateData.Properties.LayoutToUse, "Help")

}

func TestNewAPLAskResponse(t *testing.T) {

	FileToRead = "../apl_template_export.json"
	customDisplayData := CustomDataToDisplay{
		ItemsListContent: make([]string, 3),
	}
	customDisplayData.ItemsListContent[0] = "Test Content Filler"
	sessAttrData := make(map[string]interface{})

	response := NewAPLAskResponse("TextTitle1", "<speak>TextToSpeak</speak>", "<speak>TextToReprompt</speak>", "TextToDisplay", false, &sessAttrData, "Home", &customDisplayData)
	assert.Equal(t, "Alexa.Presentation.APL.RenderDocument", response.Body.Directives[0].Type, "Alexa.Presentation.APL.RenderDocument")
	assert.Contains(t, response.Body.Reprompt.OutputSpeech.SSML, "<speak>TextToReprompt</speak>")
	assert.Equal(t, "Home", response.Body.Directives[0].DataSources.TemplateData.Properties.LayoutToUse, "Home")
	response = NewAPLAskResponse("TextTitle1", "<speak>TextToSpeak</speak>", "<speak>TextToReprompt</speak>", "TextToDisplay", false, &sessAttrData, "ItemsList", &customDisplayData)
	assert.Equal(t, "Alexa.Presentation.APL.RenderDocument", response.Body.Directives[0].Type, "Alexa.Presentation.APL.RenderDocument")
	assert.Contains(t, response.Body.Reprompt.OutputSpeech.SSML, "<speak>TextToReprompt</speak>")
	assert.Equal(t, "ItemsList", response.Body.Directives[0].DataSources.TemplateData.Properties.LayoutToUse, "ItemsList")
	response = NewAPLAskResponse("TextTitle1", "<speak>TextToSpeak</speak>", "<speak>TextToReprompt</speak>", "TextToDisplay", false, &sessAttrData, "Help", &customDisplayData)
	assert.Equal(t, "Alexa.Presentation.APL.RenderDocument", response.Body.Directives[0].Type, "Alexa.Presentation.APL.RenderDocument")
	assert.Contains(t, response.Body.Reprompt.OutputSpeech.SSML, "<speak>TextToReprompt</speak>")
	assert.Equal(t, "Help", response.Body.Directives[0].DataSources.TemplateData.Properties.LayoutToUse, "Help")

}

func TestNewAPLTellResponseHandlesAPLError(t *testing.T) {

	FileToRead = "../apl_template_export.json"
	customDisplayData := CustomDataToDisplay{
		ItemsListContent: make([]string, 3),
	}
	customDisplayData.ItemsListContent[0] = "Test Content Filler"
	sessAttrData := make(map[string]interface{})
	FetchAPL = func() (*APLDocumentAndData, error) {
		return nil, errors.New("Mock APL Fetch Error")
	}
	response := NewAPLTellResponse("TextTitle1", "<speak>TextToSpeak</speak>", "TextToDisplay", false, &sessAttrData, "Home", &customDisplayData)
	assert.Empty(t, response.Body.Directives, "Should not contain a Directive array item for APL Document")
}

func TestNewAPLAskResponseHandlesAPLError(t *testing.T) {

	FileToRead = "../apl_template_export.json"
	customDisplayData := CustomDataToDisplay{
		ItemsListContent: make([]string, 3),
	}
	customDisplayData.ItemsListContent[0] = "Test Content Filler"
	sessAttrData := make(map[string]interface{})
	FetchAPL = func() (*APLDocumentAndData, error) {
		return nil, errors.New("Mock APL Fetch Error")
	}
	response := NewAPLAskResponse("TextTitle1", "<speak>TextToSpeak</speak>", "<speak>TextToReprompt</speak>", "TextToDisplay", false, &sessAttrData, "Home", &customDisplayData)
	assert.Empty(t, response.Body.Directives, "Should not contain a Directive array item for APL Document")

}

func TestParseString(t *testing.T) {
	response := ParseString("This & That")
	assert.Equal(t, "this and that", response)
}
