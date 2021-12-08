package ssml

// ssml builder implementation taken from https://www.thepolyglotdeveloper.com/2018/10/build-alexa-skill-golang-aws-lambda/
// Then rewritten to be simpler and use strings.Builder

import (
	"strings"

	"github.com/arienmalec/alexa-go"
)

const (
	WeakBreakStrength   = "weak"
	StrongBreakStrength = "strong"
	OneSecondBreak      = "1000"
)

func NewSSMLResponse(title string, text string) alexa.Response {
	r := alexa.Response{
		Version: "1.0",
		Body: alexa.ResBody{
			OutputSpeech: &alexa.Payload{
				Type: "SSML",
				SSML: text,
			},
			ShouldEndSession: true,
		},
	}
	return r
}

func NewSSMLBuilder() *SSMLBuilder {
	builder := &SSMLBuilder{}
	builder.sb.WriteString("<speak>")
	return builder
}

type SSMLBuilder struct {
	sb strings.Builder
}

func (builder *SSMLBuilder) Say(text string) *SSMLBuilder {
	builder.sb.WriteString(text)
	builder.sb.WriteString(" ")
	return builder
}

func (builder *SSMLBuilder) PauseTime(pause string) *SSMLBuilder {
	builder.sb.WriteString("<break time='")
	builder.sb.WriteString(pause)
	builder.sb.WriteString("ms'/> ")
	return builder
}

func (builder *SSMLBuilder) PauseStrength(pause string) *SSMLBuilder {
	builder.sb.WriteString("<break strength='")
	builder.sb.WriteString(pause)
	builder.sb.WriteString("'/> ")
	return builder
}

func (builder *SSMLBuilder) Build() string {
	builder.sb.WriteString("</speak>")
	return builder.sb.String()
}

func (builder *SSMLBuilder) ToResponse(title string) alexa.Response {
	return NewSSMLResponse(title, builder.Build())
}
