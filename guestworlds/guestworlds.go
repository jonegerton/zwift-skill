package guestworlds

import (
	"time"
	"zwiftaideskill/ssml"

	"github.com/arienmalec/alexa-go"
)

type guestItem struct {
	StartTime time.Time
	WorldOne  string
	WorldTwo  string
}

//Define a constant for each of the guest worlds possible
const (
	makuriIslands = "Makuri Islands"
	newYork       = "New York"
	france        = "France"
	innsbruck     = "Innsbruck"
	richmond      = "Richmond"
	london        = "London"
	paris         = "Paris"
	yorkshire     = "Yorkshire"
)

var (
	allGuestWorlds = map[string]struct{}{
		makuriIslands: {},
		newYork:       {},
		france:        {},
		innsbruck:     {},
		richmond:      {},
		london:        {},
		paris:         {},
		yorkshire:     {},
	}

	// Export a time now func reference for test overriding
	timeNow = time.Now
)

// Switch calendar switches at UTC-4 (Eastern Daylight Time)
var calendarLocation *time.Location
var calendar []guestItem

func init() {
	calendarLocation, _ = time.LoadLocation("Etc/GMT+4")
	calendar = GetCalendar()
}

//Makes a date, offset into the correct timezone for guest world change overs
func makeDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, calendarLocation)
}

func HandleGuestWorldsNowIntent() alexa.Response {

	now := timeNow()

	var guest guestItem
	for _, item := range calendar {

		if item.StartTime.After(now) {
			break
		}

		guest = item
	}

	return ssml.NewSSMLBuilder().
		Say("The current guest worlds on Zwift").
		PauseStrength(ssml.WeakBreakStrength).
		Say("are").
		Say(guest.WorldOne).
		PauseStrength(ssml.WeakBreakStrength).
		Say("and").
		Say(guest.WorldTwo).
		ToResponse("Guest Worlds Now")
}

func HandleGuestWorldsNextIntent() alexa.Response {

	now := timeNow()

	next := -1 //Enable check whether we found any items after today
	for i := len(calendar) - 1; i >= 0; i-- {
		if calendar[i].StartTime.Before(now) || calendar[i].StartTime.Equal(now) {
			break
		}
		next = i
	}

	if next >= 0 {
		guest := calendar[next]

		return ssml.NewSSMLBuilder().
			Say("The next guest worlds on Zwift").
			PauseStrength(ssml.WeakBreakStrength).
			Say("will be").
			Say(guest.WorldOne).
			PauseStrength(ssml.WeakBreakStrength).
			Say("and").
			Say(guest.WorldTwo).
			Say(".").
			PauseStrength(ssml.StrongBreakStrength).
			Say("They will be available from").
			Say(guest.StartTime.Format("January 2")).
			ToResponse("Guest Worlds Next")
	}

	// No next track is currently available (Zwift can leave calendar releases very late)
	return ssml.NewSSMLBuilder().
		Say("The next guest worlds are not currently available").
		ToResponse("Guest Worlds Next")
}

func HandleGuestWorldsDateIntent(request alexa.Request) alexa.Response {

	// We could check the slot resolutions fields, but easier just to check
	// the value against a whitelist.
	dateSlot, ok := request.Body.Intent.Slots["Date"]
	if !ok {
		//Not sure if this can happen
		return ssml.NewSSMLBuilder().
			Say("I didn't understand which date you asked about").
			PauseStrength(ssml.StrongBreakStrength).
			Say("Please try asking again.").
			ToResponse("Guest World Date")
	}
	date, err := time.Parse("2006-01-02", dateSlot.Value)
	if err != nil {

		//Value requested wasn't recongnised as a guest world.
		return ssml.NewSSMLBuilder().
			Say("I didn't understand which date you asked about").
			PauseStrength(ssml.StrongBreakStrength).
			Say("Please try asking again.").
			ToResponse("Guest World Date")
	}

	// Shift dates to midday so are after 4am changeover
	// This means we get the more obvious result
	date = date.Add(12 * time.Hour)

	var guest guestItem
	for _, item := range calendar {

		if item.StartTime.After(date) {
			break
		}

		guest = item
	}

	//Check we resolved an item
	if guest.WorldOne == "" {
		return ssml.NewSSMLBuilder().
			Say("I couldn't find the guest worlds available on ").
			Say(date.Format("January 2")).
			ToResponse("Guest Worlds Date")
	}

	return ssml.NewSSMLBuilder().
		Say("The guest worlds on Zwift").
		PauseStrength(ssml.WeakBreakStrength).
		Say("on").
		Say(date.Format("January 2")).
		PauseStrength(ssml.WeakBreakStrength).
		Say("are").
		Say(guest.WorldOne).
		PauseStrength(ssml.WeakBreakStrength).
		Say("and").
		Say(guest.WorldTwo).
		ToResponse("Guest Worlds Date")
}

func HandleWhenGuestWorldIntent(request alexa.Request) alexa.Response {

	// We could check the slot resolutions fields, but easier just to check
	// the value against a whitelist.
	whenGuestSlot, ok := request.Body.Intent.Slots["GuestWorld"]
	if !ok {
		//Not sure if this can happen
		return ssml.NewSSMLBuilder().
			Say("I'm not sure which guest world you asked about").
			PauseStrength(ssml.StrongBreakStrength).
			Say("Please try asking again.").
			ToResponse("Guest World When")
	}
	whenGuest := whenGuestSlot.Value
	if _, ok := allGuestWorlds[whenGuest]; !ok {

		//Value requested wasn't recongnised as a guest world.
		return ssml.NewSSMLBuilder().
			Say("I'm not sure which guest world you asked about").
			PauseStrength(ssml.StrongBreakStrength).
			Say("Please try asking again.").
			ToResponse("Guest World When")
	}

	now := timeNow()
	next := -1 //Enable check whether we found any items after today
	for i := len(calendar) - 1; i >= 0; i-- {
		item := calendar[i]
		if item.WorldOne == whenGuest || item.WorldTwo == whenGuest {
			next = i
		}
		if item.StartTime.Before(now) || item.StartTime.Equal(now) {
			break
		}
	}

	if next >= 0 {
		guest := calendar[next]

		return ssml.NewSSMLBuilder().
			Say(whenGuest).
			Say("will next be available on").
			Say(guest.StartTime.Format("January 2")).
			ToResponse("Guest World When")

	}

	// Guest world is not currently available (Zwift can leave calendar releases very late)
	return ssml.NewSSMLBuilder().
		Say(whenGuest).
		Say("has no dates available in the zwift guest world calendar.").
		ToResponse("Guest World When")
}
