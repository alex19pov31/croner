package croner

import (
	"fmt"
	"testing"
	"time"
)

func TestSimpleForm(t *testing.T) {
	chTime, err := time.Parse(time.RFC3339, "2019-01-02T21:22:05Z")
	if err != nil {
		t.Error(err)
	}

	if !NewCronTimer("22 21 2 1 *").Check(chTime) {
		t.Error("Invalid check")
	}

	if !NewCronTimer("22 21 2 1 3").Check(chTime) {
		t.Error("Invalid check")
	}

	if NewCronTimer("22 21 2 1 4").Check(chTime) {
		t.Error("Invalid check")
	}

	if NewCronTimer("22 21 2 2 *").Check(chTime) {
		t.Error("Invalid check")
	}

	if NewCronTimer("22 21 3 1 *").Check(chTime) {
		t.Error("Invalid check")
	}

	if NewCronTimer("22 22 2 1 *").Check(chTime) {
		t.Error("Invalid check")
	}

	if NewCronTimer("24 21 2 1 *").Check(chTime) {
		t.Error("Invalid check")
	}
}

func TestMultiplyMinutes(t *testing.T) {
	cronTimer := NewCronTimer("*/2,3,7 21 2 1 *")
	if cronTimer.Check(getTime("2019-01-02T21:01:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:02:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:03:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:04:05Z", t)) {
		t.Error("Invalid check")
	}

	if cronTimer.Check(getTime("2019-01-02T21:05:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:06:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:07:05Z", t)) {
		t.Error("Invalid check")
	}
}
func TestMultiplyHours(t *testing.T) {
	cronTimer := NewCronTimer("* */2,3,7 2 1 *")
	if cronTimer.Check(getTime("2019-01-02T01:01:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T02:02:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T03:03:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T04:04:05Z", t)) {
		t.Error("Invalid check")
	}

	if cronTimer.Check(getTime("2019-01-02T05:05:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T06:06:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T07:07:05Z", t)) {
		t.Error("Invalid check")
	}
}

func TestMultiplyDays(t *testing.T) {
	cronTimer := NewCronTimer("* * */2,3,7 1 *")
	if cronTimer.Check(getTime("2019-01-01T21:01:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:02:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-03T21:03:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-04T21:04:05Z", t)) {
		t.Error("Invalid check")
	}

	if cronTimer.Check(getTime("2019-01-05T21:05:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-06T21:06:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-07T21:07:05Z", t)) {
		t.Error("Invalid check")
	}
}

func TestLimit(t *testing.T) {
	cronTimer := NewCronTimer("15-18 21 2 1 *")

	if cronTimer.Check(getTime("2019-01-02T21:01:05Z", t)) {
		t.Error("Invalid check")
	}

	if cronTimer.Check(getTime("2019-01-02T21:04:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:15:05Z", t)) {
		t.Error("Invalid check")
	}
	if !cronTimer.Check(getTime("2019-01-02T21:16:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:17:05Z", t)) {
		t.Error("Invalid check")
	}

	if !cronTimer.Check(getTime("2019-01-02T21:18:05Z", t)) {
		t.Error("Invalid check")
	}

	data := map[string]int{"string": 1}
	NewCronTimer("15-18 21 2 1 *").Start(func(data interface{}) {
		d, ok := data.(map[string]int)
		if !ok {
			return
		}

		fmt.Println(d)
	}, data)
}

func TestSlower(t *testing.T) {
	cronTimer := NewCronTimer("15-18,*/3,11 21 2 1 *")
	if !cronTimer.Check(getTime("2019-01-02T21:06:05Z", t)) {
		t.Error("Invalid check")
	}
}

func getTime(strTime string, t *testing.T) time.Time {
	chTime, err := time.Parse(time.RFC3339, strTime)
	if err != nil {
		t.Error(err)
	}

	return chTime
}
