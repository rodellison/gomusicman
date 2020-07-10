package alexa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanLoadAPL(t *testing.T) {

	//For test purposes, the file to read is actually sitting in the main directory, not the package so fixing for
	//test accordingly..
	FileToRead = "../apl_template_export.json"
	response, _ := FetchAPL()
	assert.Equal(t, "APL", response.APLDocument.Type, "Can create APL structure")
	assert.Equal(t, "1.0", response.APLDocument.Version, "Can create APL structure")
}

func TestCanHandleAPLFileNotFound(t *testing.T) {

	//For test purposes, the file to read is actually sitting in the main directory, not the package so fixing for
	//test accordingly..
	FileToRead = "apl_template_export" //bad file name
	_, err := FetchAPL()
	assert.Error(t, err)
}
