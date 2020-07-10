package alexa

type SSML struct {
	text  string
	pause string
}

type SSMLBuilder struct {
	SSML []SSML
}

//SSMLBuilder function for adding new SSML content that will be 'spoken' back to the the user in the Response
func (builder *SSMLBuilder) Say(text string) {
	//NOTE!! ParseString result is a Lowercase string
	text = ParseString(text)
	builder.SSML = append(builder.SSML, SSML{text: text})
}

//SSMLBuilder function for adding pause/delay content that will be used in creating audible delays within the'spoken' provided
//back to the the user in the Response
func (builder *SSMLBuilder) Pause(pause string) {
	builder.SSML = append(builder.SSML, SSML{pause: pause})
}

//SSMLBuilder function for producing the overall '<speak>' SSML content, including all pauses, etc and providing back
//to caller as a complete text string.
func (builder *SSMLBuilder) Build() string {
	var response string
	for index, ssml := range builder.SSML {
		if ssml.text != "" {
			response += ssml.text + " "
		} else if ssml.pause != "" && index != len(builder.SSML)-1 {
			response += "<break time='" + ssml.pause + "ms'/> "
		}
	}
	return "<speak>" + response + "</speak>"
}
