package guestworlds

import (
	"strings"
	"testing"
	"time"

	"github.com/arienmalec/alexa-go"
)

func TestHandleGuestWorldsNowIntent(t *testing.T) {

	//Set up dummy calendar items for test
	calendar = []guestItem{
		{StartTime: makeDate(2021, time.December, 1), WorldOne: makuriIslands, WorldTwo: newYork},
		{StartTime: makeDate(2021, time.December, 3), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 5), WorldOne: france, WorldTwo: paris},
	}

	utc, _ := time.LoadLocation("UTC")
	tcs := []struct {
		name    string
		dateNow time.Time
		expect  []string
	}{
		{
			name:    "Date before 4am GMT yields correct response",
			dateNow: time.Date(2021, time.December, 3, 3, 0, 0, 0, utc),
			expect:  []string{makuriIslands, newYork},
		},
		{
			name:    "Date at 4am GMT yields correct response",
			dateNow: time.Date(2021, time.December, 3, 4, 0, 0, 0, utc),
			expect:  []string{innsbruck, richmond},
		},
		{
			name:    "Date after 4am GMT yields correct response",
			dateNow: time.Date(2021, time.December, 3, 12, 0, 0, 0, utc),
			expect:  []string{innsbruck, richmond},
		},
		{
			name:    "Date after all calendar items yields correct response",
			dateNow: time.Date(2021, time.December, 6, 12, 0, 0, 0, utc),
			expect:  []string{france, paris},
		},
	}

	for _, tc := range tcs {

		timeNow = func() time.Time {
			return tc.dateNow
		}

		t.Run(tc.name, func(t *testing.T) {

			response := HandleGuestWorldsNowIntent()

			got := response.Body.OutputSpeech.SSML

			if !containsAll(got, tc.expect) {
				t.Errorf("Expected response '%s', got '%s'", tc.expect, got)
			}

		})
	}
}

func TestHandleGuestWorldsNextIntent(t *testing.T) {

	//Set up dummy calendar items for test
	calendar = []guestItem{
		{StartTime: makeDate(2021, time.December, 1), WorldOne: makuriIslands, WorldTwo: newYork},
		{StartTime: makeDate(2021, time.December, 3), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 5), WorldOne: france, WorldTwo: paris},
	}

	utc, _ := time.LoadLocation("UTC")
	tcs := []struct {
		name    string
		dateNow time.Time
		expect  []string
	}{
		{
			name:    "Date before 4am GMT yields correct response",
			dateNow: time.Date(2021, time.December, 3, 3, 0, 0, 0, utc),
			expect:  []string{innsbruck, richmond, "December 3"},
		},
		{
			name:    "Date at 4am GMT yields correct response",
			dateNow: time.Date(2021, time.December, 3, 4, 0, 0, 0, utc),
			expect:  []string{france, paris, "December 5"},
		},
		{
			name:    "Date after 4am GMT yields correct response",
			dateNow: time.Date(2021, time.December, 3, 12, 0, 0, 0, utc),
			expect:  []string{france, paris, "December 5"},
		},
		{
			name:    "Date after all calendar items yields correct response",
			dateNow: time.Date(2021, time.December, 6, 12, 0, 0, 0, utc),
			expect:  []string{"not currently available"},
		},
	}

	for _, tc := range tcs {

		timeNow = func() time.Time {
			return tc.dateNow
		}

		t.Run(tc.name, func(t *testing.T) {

			response := HandleGuestWorldsNextIntent()

			got := response.Body.OutputSpeech.SSML

			if !containsAll(got, tc.expect) {
				t.Errorf("Expected response '%v', got '%s'", tc.expect, got)
			}

		})
	}
}

func TestHandleGuestWorldsDateIntent(t *testing.T) {

	//Set up dummy calendar items for test
	calendar = []guestItem{
		{StartTime: makeDate(2021, time.December, 1), WorldOne: makuriIslands, WorldTwo: newYork},
		{StartTime: makeDate(2021, time.December, 3), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 5), WorldOne: france, WorldTwo: paris},
	}

	utc, _ := time.LoadLocation("UTC")
	tcs := []struct {
		name      string
		guestDate string
		expect    []string
	}{
		{
			name:      "Date before any available yields correct response",
			guestDate: time.Date(2021, time.November, 29, 0, 0, 0, 0, utc).Format("2006-01-02"),
			expect:    []string{"couldn't find"},
		},
		{
			name:      "Date mid calendar yields correct response",
			guestDate: time.Date(2021, time.December, 3, 0, 0, 0, 0, utc).Format("2006-01-02"),
			expect:    []string{innsbruck, richmond},
		},
		{
			name:      "Date after calendar yields correct response",
			guestDate: time.Date(2021, time.December, 5, 0, 0, 0, 0, utc).Format("2006-01-02"),
			expect:    []string{france, paris},
		},
		{
			name:      "Request for world with no date yields correct response",
			guestDate: "",
			expect:    []string{"didn't understand"},
		},
		{
			name:      "Request for incorrect date format yields correct response",
			guestDate: "2021-12-XX",
			expect:    []string{"didn't understand"},
		},
	}

	for _, tc := range tcs {

		t.Run(tc.name, func(t *testing.T) {

			request := alexa.Request{
				Body: alexa.ReqBody{
					Intent: alexa.Intent{
						Slots: make(map[string]alexa.Slot),
					},
				},
			}

			if tc.guestDate != "" {
				request.Body.Intent.Slots["Date"] = alexa.Slot{
					Name:  "Date",
					Value: tc.guestDate,
				}
			}

			response := HandleGuestWorldsDateIntent(request)

			got := response.Body.OutputSpeech.SSML

			if !containsAll(got, tc.expect) {
				t.Errorf("Expected response '%v', got '%s'", tc.expect, got)
			}

		})
	}
}

func TestHandleWhenGuestWorldIntent(t *testing.T) {

	//Set up dummy calendar items for test
	calendar = []guestItem{
		{StartTime: makeDate(2021, time.December, 1), WorldOne: makuriIslands, WorldTwo: newYork},
		{StartTime: makeDate(2021, time.December, 3), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 5), WorldOne: france, WorldTwo: paris},
	}

	utc, _ := time.LoadLocation("UTC")
	tcs := []struct {
		name      string
		whenGuest string
		dateNow   time.Time
		expect    []string
	}{
		{
			name:      "Date before 4am GMT yields correct response",
			whenGuest: richmond,
			dateNow:   time.Date(2021, time.December, 3, 3, 0, 0, 0, utc),
			expect:    []string{richmond, "December 3"},
		},
		{
			name:      "Date at 4am GMT yields correct response",
			whenGuest: paris,
			dateNow:   time.Date(2021, time.December, 3, 4, 0, 0, 0, utc),
			expect:    []string{paris, "December 5"},
		},
		{
			name:      "Today's guest world and date after 4am GMT yields correct response",
			whenGuest: innsbruck,
			dateNow:   time.Date(2021, time.December, 3, 5, 0, 0, 0, utc),
			expect:    []string{innsbruck, "December 3"},
		},
		{
			name:      "Request for world with no date yields correct response",
			whenGuest: makuriIslands,
			dateNow:   time.Date(2021, time.December, 3, 12, 0, 0, 0, utc),
			expect:    []string{"has no dates"},
		},
		{
			name:      "Request for incorrect guest world yields correct response",
			whenGuest: "foobar",
			dateNow:   time.Date(2021, time.December, 6, 12, 0, 0, 0, utc),
			expect:    []string{"not sure which"},
		},
	}

	for _, tc := range tcs {

		timeNow = func() time.Time {
			return tc.dateNow
		}

		t.Run(tc.name, func(t *testing.T) {

			request := alexa.Request{
				Body: alexa.ReqBody{
					Intent: alexa.Intent{
						Slots: map[string]alexa.Slot{
							"GuestWorld": {
								Name:  "GuestWorld",
								Value: tc.whenGuest,
							},
						},
					},
				},
			}
			response := HandleWhenGuestWorldIntent(request)

			got := response.Body.OutputSpeech.SSML

			if !containsAll(got, tc.expect) {
				t.Errorf("Expected response '%v', got '%s'", tc.expect, got)
			}

		})
	}
}

func containsAll(s string, contains []string) bool {

	for _, c := range contains {
		if !strings.Contains(s, c) {
			return false
		}
	}

	return true
}
