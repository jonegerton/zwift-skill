package guestworlds

import "time"

func GetCalendar() []guestItem {
	return []guestItem{
		{StartTime: makeDate(2021, time.December, 1), WorldOne: makuriIslands, WorldTwo: newYork},
		{StartTime: makeDate(2021, time.December, 3), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 5), WorldOne: france, WorldTwo: paris},
		{StartTime: makeDate(2021, time.December, 7), WorldOne: london, WorldTwo: yorkshire},
		{StartTime: makeDate(2021, time.December, 10), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 12), WorldOne: yorkshire, WorldTwo: innsbruck},
		{StartTime: makeDate(2021, time.December, 14), WorldOne: makuriIslands, WorldTwo: newYork},
		{StartTime: makeDate(2021, time.December, 17), WorldOne: france, WorldTwo: paris},
		{StartTime: makeDate(2021, time.December, 19), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 21), WorldOne: london, WorldTwo: yorkshire},
		{StartTime: makeDate(2021, time.December, 24), WorldOne: makuriIslands, WorldTwo: newYork},
		{StartTime: makeDate(2021, time.December, 25), WorldOne: yorkshire, WorldTwo: innsbruck},
		{StartTime: makeDate(2021, time.December, 26), WorldOne: france, WorldTwo: paris},
		{StartTime: makeDate(2021, time.December, 27), WorldOne: innsbruck, WorldTwo: richmond},
		{StartTime: makeDate(2021, time.December, 28), WorldOne: makuriIslands, WorldTwo: newYork},
	}
}
