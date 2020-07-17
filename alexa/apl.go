package alexa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	//FetchAPL variable is defined as a func() that returns APLDocumentAndData structure. It is a variable so that it can be overridden for testing purposes
	FetchAPL   func() (*APLDocumentAndData, error)
	FileToRead string
)

func init() {
	//FetchAPL is initialized to the CreateAPLDocAndData function, but can be overridden or mocked out for testing purposes as needed.
	FetchAPL = CreateAPLDocAndData
	//FileToRead is initialized by way of an Environmental variable, but can be overridden as necessary
	FileToRead = os.Getenv("AplTemplate")
}

type APLDocumentAndData struct {
	APLDocument    APLDocument    `json:"document"`
	APLDataSources APLDataSources `json:"datasources"`
	APLSources     interface{}    `json:"sources,omitempty"`
}

//The types defined as []interface{} or just interface[] can stay as is, there's no need to break them down to
//struct level as they shouldn't be modified by code after loading
//They are broken to this level just to facilitate some file load and UNMarshal testing
type APLDocument struct {
	Type         string        `json:"type,omitempty"`
	Version      string        `json:"version,omitempty"`
	Theme        string        `json:"theme,omitempty"`
	Import       []interface{} `json:"import,omitempty"`
	Resources    []interface{} `json:"resources,omitempty"`
	Styles       interface{}   `json:"styles,omitempty"`
	Layouts      interface{}   `json:"layouts,omitempty"`
	MainTemplate interface{}   `json:"mainTemplate,omitempty"`
}

//The DataSources sections and types below needs to be adjusted to match the apl_template_export.json file as the properties and variables for the APL document
//are completely customized and unique for each skill.
//We need the properties to be accessible so as to dynamically populate them when responding to user requests
type APLDataSources struct {
	TemplateData TemplateData `json:"musicManTemplateData,omitempty"` //NOTE!, make sure to use the correct json: value if different than templateData!
}

type TemplateData struct {
	Type         string            `json:"type"`
	ObjectID     string            `json:"objectId"`
	Properties   APLDataProperties `json:"properties,omitempty"`
	Transformers []interface{}     `json:"transformers,omitempty"`
}

//These properties need to match up to those in the DataSources section of the json apl_template_export.json file
type APLDataProperties struct {
	Title                    string   `json:"Title"`
	LayoutToUse              string   `json:"LayoutToUse"`
	HeadingText              string   `json:"HeadingText"`
	Locale                   string   `json:"Locale"`
	HintString               string   `json:"HintString"`
	EventImageUrl            string   `json:"EventImageUrl"`
	BackgroundImageUrl       string   `json:"BackgroundImageUrl"`
	LogoUrl                  string   `json:"LogoUrl"`
	GeneralSquareImageUrl    string   `json:"GeneralSquareImageUrl"`
	SongkickLogoUrl          string   `json:"SongkickLogoUrl"`
	EventText                []string `json:"EventText"`
	PhotoAttribution         string `json:"PhotoAttribution"`
	BackgroundImages         []string `json:"BackgroundImages"`
}

//This struct will be used to define a container type for passing custom display data to the NewTell/Ask..Response functions
//It should be customized as needed to support the data needs of the dynamically updated Datasource APL properties
type CustomDataToDisplay struct {
	ItemsListContent  []string
	ArtistVenueImgURL string
}

//CreateAPLDocAndData is a function that loads json input that contains the APL template for the Skill. Once loaded, the content
//is unmarshalled to an overall APLDocAndData structure, where it can then be manipulated in code before inserting into the
//Response structure which is passed back to the Alexa device.
func CreateAPLDocAndData() (*APLDocumentAndData, error) {

	aplDoc := &APLDocumentAndData{
		APLDocument:    APLDocument{},
		APLDataSources: APLDataSources{},
	}

	//apl_template_export.json is the name given to the file downloaded from the Alexa Developer console (Display tool)
	//fmt.Println("Opening APLTemplate JSON file: ", FileToRead)
	file, err := ioutil.ReadFile(FileToRead)
	if err != nil {
		fmt.Println("Error reading APL json file: ", err.Error())
		return nil, err
	} else {
		errUnMarshal := json.Unmarshal([]byte(file), &aplDoc)
		if errUnMarshal != nil {
			fmt.Println("Error Unmarshalling APL json file: ", err.Error())
			return nil, err
		}
	}

	return aplDoc, nil

}
