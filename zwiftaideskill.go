package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"zwiftaideskill/guestworlds"
	"zwiftaideskill/ssml"

	"github.com/arienmalec/alexa-go"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func Handler(request alexa.Request) (alexa.Response, error) {
	var response alexa.Response
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v, %s", r, string(debug.Stack()))
		}
		if err != nil {
			response = alexa.NewSimpleResponse("Oops", "Zwift Aide had a problem.")
			//Need to log and return boiler plate sorry response
		}
	}()

	switch request.Body.Intent.Name {
	case "GuestWorldsNowIntent":
		response = guestworlds.HandleGuestWorldsNowIntent()
	case "GuestWorldsNextIntent":
		response = guestworlds.HandleGuestWorldsNextIntent()
	case "GuestWorldsDateIntent":
		response = guestworlds.HandleGuestWorldsDateIntent(request)
	case "WhenGuestWorldIntent":
		response = guestworlds.HandleWhenGuestWorldIntent(request)
	case alexa.HelpIntent:
		response = handleHelpIntent(request)
	default:
		response = handleHelpIntent(request)
	}

	return response, err
}

var helpResponse alexa.Response
var helpResponseSet bool
var helpResponseLock sync.Mutex

func handleHelpIntent(request alexa.Request) alexa.Response {

	if !helpResponseSet {
		helpResponseLock.Lock()
		defer helpResponseLock.Unlock()

		if !helpResponseSet {
			helpResponse = ssml.NewSSMLBuilder().
				Say("Here are some of the things you can ask:").
				PauseStrength(ssml.StrongBreakStrength).
				Say("What are the guest worlds now?").
				PauseStrength(ssml.StrongBreakStrength).
				Say("What are the next guest worlds?").
				PauseStrength(ssml.StrongBreakStrength).
				Say("Which worlds are available on December 10th?").
				PauseStrength(ssml.StrongBreakStrength).
				Say("When can I ride Richmond?").
				ToResponse("Aide for Zwift help")
			helpResponseSet = true
		}
	}
	return helpResponse
}
