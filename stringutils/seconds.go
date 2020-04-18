package stringutils

import (
	"math"
	"strconv"
)

/* Credits: https://www.socketloop.com/tutorials/golang-convert-seconds-to-human-readable-time-format-example */

func Pluralize(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func SecondsToHuman(input int64) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if years > 0 {
		result = Pluralize(int(years), "year") + Pluralize(int(months), "month") + Pluralize(int(weeks), "week") + Pluralize(int(days), "day") + Pluralize(int(hours), "hour") + Pluralize(int(minutes), "minute") + Pluralize(int(seconds), "second")
	} else if months > 0 {
		result = Pluralize(int(months), "month") + Pluralize(int(weeks), "week") + Pluralize(int(days), "day") + Pluralize(int(hours), "hour") + Pluralize(int(minutes), "minute") + Pluralize(int(seconds), "second")
	} else if weeks > 0 {
		result = Pluralize(int(weeks), "week") + Pluralize(int(days), "day") + Pluralize(int(hours), "hour") + Pluralize(int(minutes), "minute") + Pluralize(int(seconds), "second")
	} else if days > 0 {
		result = Pluralize(int(days), "day") + Pluralize(int(hours), "hour") + Pluralize(int(minutes), "minute") + Pluralize(int(seconds), "second")
	} else if hours > 0 {
		result = Pluralize(int(hours), "hour") + Pluralize(int(minutes), "minute") + Pluralize(int(seconds), "second")
	} else if minutes > 0 {
		result = Pluralize(int(minutes), "minute") + Pluralize(int(seconds), "second")
	} else {
		result = Pluralize(int(seconds), "second")
	}

	return
}
