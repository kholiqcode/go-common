package common_utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddDay(t *testing.T) {

	now := time.Now()
	after2Days := AddDay(now, 2)
	diff1 := after2Days.Sub(now)
	assert.Equal(t, float64(48), diff1.Hours(), "TestAddDay")

	before2Days := AddDay(now, -2)
	diff2 := before2Days.Sub(now)
	assert.Equal(t, float64(-48), diff2.Hours(), "TestAddDay")
}

func TestAddHour(t *testing.T) {

	now := time.Now()
	after2Hours := AddHour(now, 2)
	diff1 := after2Hours.Sub(now)
	assert.Equal(t, float64(2), diff1.Hours(), "TestAddHour")

	before2Hours := AddHour(now, -2)
	diff2 := before2Hours.Sub(now)
	assert.Equal(t, float64(-2), diff2.Hours(), "TestAddHour")
}

func TestAddMinute(t *testing.T) {

	now := time.Now()
	after2Minutes := AddMinute(now, 2)
	diff1 := after2Minutes.Sub(now)
	assert.Equal(t, float64(2), diff1.Minutes(), "TestAddMinute")

	before2Minutes := AddMinute(now, -2)
	diff2 := before2Minutes.Sub(now)
	assert.Equal(t, float64(-2), diff2.Minutes(), "TestAddMinute")
}

func TestGetNowDate(t *testing.T) {
	expected := time.Now().Format("2006-01-02")
	assert.Equal(t, expected, GetNowDate(), "TestGetNowDate")
}

func TestGetNotTime(t *testing.T) {
	expected := time.Now().Format("15:04:05")
	assert.Equal(t, expected, GetNowTime(), "TestGetNowTime")
}

func TestGetNowDateTime(t *testing.T) {
	expected := time.Now().Format("2006-01-02 15:04:05")
	assert.Equal(t, expected, GetNowDateTime(), "TestGetNowDateTime")
}

func TestFormatTimeToStr(t *testing.T) {

	datetime, _ := time.Parse("2006-01-02 15:04:05", "2021-01-02 16:04:08")
	cases := []string{
		"yyyy-mm-dd hh:mm:ss", "yyyy-mm-dd",
		"dd-mm-yy hh:mm:ss", "yyyy/mm/dd hh:mm:ss",
		"hh:mm:ss", "yyyy/mm"}

	expected := []string{
		"2021-01-02 16:04:08", "2021-01-02",
		"02-01-21 16:04:08", "2021/01/02 16:04:08",
		"16:04:08", "2021/01"}

	for i := 0; i < len(cases); i++ {
		actual := FormatTimeToStr(datetime, cases[i])
		assert.Equal(t, expected[i], actual, "TestFormatTimeToStr")

	}
}

func TestFormatStrToTime(t *testing.T) {

	formats := []string{
		"2006-01-02 15:04:05", "2006-01-02",
		"02-01-06 15:04:05", "2006/01/02 15:04:05",
		"2006/01"}
	cases := []string{
		"yyyy-mm-dd hh:mm:ss", "yyyy-mm-dd",
		"dd-mm-yy hh:mm:ss", "yyyy/mm/dd hh:mm:ss",
		"yyyy/mm"}

	datetimeStr := []string{
		"2021-01-02 16:04:08", "2021-01-02",
		"02-01-21 16:04:08", "2021/01/02 16:04:08",
		"2021/01"}

	for i := 0; i < len(cases); i++ {
		actual, err := FormatStrToTime(datetimeStr[i], cases[i])
		if err != nil {
			t.Fatal(err)
		}
		expected, _ := time.Parse(formats[i], datetimeStr[i])
		assert.Equal(t, expected, actual, "TestFormatStrToTime")
	}
}

func TestBeginOfMinute(t *testing.T) {

	expected := time.Date(2022, 2, 15, 15, 48, 0, 0, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := BeginOfMinute(td)

	assert.Equal(t, expected, actual, "TestBeginOfMinute")
}

func TestEndOfMinute(t *testing.T) {

	expected := time.Date(2022, 2, 15, 15, 48, 59, 999999999, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := EndOfMinute(td)

	assert.Equal(t, expected, actual, "TestEndOfMinute")
}

func TestBeginOfHour(t *testing.T) {

	expected := time.Date(2022, 2, 15, 15, 0, 0, 0, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := BeginOfHour(td)

	assert.Equal(t, expected, actual, "TestBeginOfHour")
}

func TestEndOfHour(t *testing.T) {

	expected := time.Date(2022, 2, 15, 15, 59, 59, 999999999, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := EndOfHour(td)

	assert.Equal(t, expected, actual, "TestEndOfHour")
}

func TestBeginOfDay(t *testing.T) {

	expected := time.Date(2022, 2, 15, 0, 0, 0, 0, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := BeginOfDay(td)

	assert.Equal(t, expected, actual, "TestBeginOfDay")
}

func TestEndOfDay(t *testing.T) {

	expected := time.Date(2022, 2, 15, 23, 59, 59, 999999999, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := EndOfDay(td)

	assert.Equal(t, expected, actual, "TestEndOfDay")
}

func TestBeginOfWeek(t *testing.T) {

	expected := time.Date(2022, 2, 13, 0, 0, 0, 0, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := BeginOfWeek(td)

	assert.Equal(t, expected, actual, "TestBeginOfWeek")
}

func TestEndOfWeek(t *testing.T) {

	expected := time.Date(2022, 2, 19, 23, 59, 59, 999999999, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := EndOfWeek(td)

	assert.Equal(t, expected, actual, "TestEndOfWeek")
}

func TestBeginOfMonth(t *testing.T) {

	expected := time.Date(2022, 2, 1, 0, 0, 0, 0, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := BeginOfMonth(td)

	assert.Equal(t, expected, actual, "TestBeginOfMonth")
}

func TestEndOfMonth(t *testing.T) {

	expected := time.Date(2022, 2, 28, 23, 59, 59, 999999999, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := EndOfMonth(td)

	assert.Equal(t, expected, actual, "TestEndOfMonth")
}

func TestBeginOfYear(t *testing.T) {

	expected := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := BeginOfYear(td)

	assert.Equal(t, expected, actual, "TestBeginOfYear")
}

func TestEndOfYear(t *testing.T) {

	expected := time.Date(2022, 12, 31, 23, 59, 59, 999999999, time.Local)
	td := time.Date(2022, 2, 15, 15, 48, 40, 112, time.Local)
	actual := EndOfYear(td)

	assert.Equal(t, expected, actual, "TestEndOfYear")
}
