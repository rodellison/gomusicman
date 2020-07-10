package alexa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanCreateSSMLContent(t *testing.T) {

	var testSSMLString SSMLBuilder
	testSSMLString.Say("Hello")
	testSSMLString.Pause("1000")
	testSSMLString.Say("World")

	assert.Len(t, testSSMLString.SSML, 3, "SSML structure contains all items added")
	assert.Contains(t, testSSMLString.SSML[1].pause, "1000", "SSML pause value is inserted")

}
func TestBuildSSML(t *testing.T) {

	var testSSMLString SSMLBuilder

	testSSMLString.Say("Hello")
	testSSMLString.Pause("1000")
	testSSMLString.Say("World")

	result := testSSMLString.Build()

	assert.Contains(t, result, "<speak>hello <break time='1000ms'/> world </speak>", "SSML build creates <speak> content")

}
