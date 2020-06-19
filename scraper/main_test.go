package main

import (
	"testing"
	"time"
)

func TestParseDateTime(t *testing.T) {
	loc, _ := time.LoadLocation("America/Montreal")

	midnight := "Dec 6th, 1994 12:05 am"
	expected := time.Date(1994, 12, 6, 0, 05, 0, 0, loc)
	t.Run("midnight", testParseDateTimeFunc(midnight, expected))

	noon := "Dec 7th, 1995 12:05 pm"
	expected = time.Date(1995, 12, 7, 12, 05, 0, 0, loc)
	t.Run("noon", testParseDateTimeFunc(noon, expected))

	standardTime := "Apr 1st, 2020 5:05 pm"
	expected = time.Date(2020, 04, 1, 17, 5, 0, 0, loc)
	t.Run("standard time", testParseDateTimeFunc(standardTime, expected))

}

func testParseDateTimeFunc(testString string, expected time.Time) func(*testing.T) {
	return func(t *testing.T) {
		result := parseDateTime(testString)
		if !result.Equal(expected) {
			t.Errorf("Expected %s, but got %s", expected, result)
		}
	}
}

func TestStrToInt(t *testing.T) {
	thousand := "1,000"
	expected := 1000
	t.Run("thousand", testStrToIntFunc(thousand, expected))

	positive := "+666"
	expected = 666
	t.Run("positive", testStrToIntFunc(positive, expected))

	negative := "-666"
	expected = -666
	t.Run("negative", testStrToIntFunc(negative, expected))

	empty := ""
	expected = 0
	t.Run("empty", testStrToIntFunc(empty, expected))

	whitespace := "   666     "
	expected = 666
	t.Run("whitespaces", testStrToIntFunc(whitespace, expected))

	decimals := "6.66"
	expected = 6
	t.Run("decimals", testStrToIntFunc(decimals, expected))

}

func testStrToIntFunc(testString string, expected int) func(*testing.T) {
	return func(t *testing.T) {
		result := strToInt(testString)
		if result != expected {
			t.Errorf("Expected %d, but got %d", expected, result)
		}
	}
}
