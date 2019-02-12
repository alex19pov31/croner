package croner

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Condition struct {
	Min          int
	Max          int
	Multiplicity int
}

// Check check limit and multiplicity
func (condition *Condition) Check(value int) bool {

	if !condition.checkMin(value) {
		return false
	}

	if !condition.checkMax(value) {
		return false
	}

	if condition.Multiplicity > 0 {
		result := condition.checkMultiplicity(value)
		return result
	}

	return true
}

func (condition *Condition) checkMin(value int) bool {
	return condition.Min == 0 || value >= condition.Min
}

func (condition *Condition) checkMax(value int) bool {
	return condition.Max == 0 || value <= condition.Max
}

func (condition *Condition) checkMultiplicity(value int) bool {
	value = value - condition.Min
	return (value % condition.Multiplicity) == 0
}

// ConditionElement element of cron record
type ConditionElement struct {
	Max        int
	Conditions []Condition
	EnumValues []int
}

// Check check cron element by cron rule
func (conditionElement *ConditionElement) Check(value int) bool {
	if value > conditionElement.Max {
		return false
	}

	if conditionElement.checkValue(value) {
		return true
	}

	for _, condition := range conditionElement.Conditions {
		if condition.Check(value) {
			return true
		}
	}

	return false
}

func (conditionElement *ConditionElement) checkValue(value int) bool {
	for _, v := range conditionElement.EnumValues {
		if value == v {
			return true
		}
	}

	return false
}

// CronTimer cron rule
type CronTimer struct {
	data             string
	minuteCondition  *ConditionElement
	hourCondition    *ConditionElement
	dayCondition     *ConditionElement
	monthCondition   *ConditionElement
	weekdayCondition *ConditionElement
	ticker           *time.Ticker
	chSuccess        chan bool
}

// Check check cron rule
func (cronTimer *CronTimer) Check(t time.Time) bool {
	elementCount := cronTimer.getCountCondition()
	if elementCount == 0 {
		return true
	}

	cronTimer.chSuccess = make(chan bool)
	go func(cronTimer *CronTimer, t time.Time) {
		go func(element *ConditionElement, value int, chSuccess chan bool) {
			if element == nil {
				return
			}
			chSuccess <- element.Check(value)
		}(cronTimer.minuteCondition, t.Minute(), cronTimer.chSuccess)

		go func(element *ConditionElement, value int, chSuccess chan bool) {
			if element == nil {
				return
			}
			chSuccess <- element.Check(value)
		}(cronTimer.hourCondition, t.Hour(), cronTimer.chSuccess)

		go func(element *ConditionElement, value int, chSuccess chan bool) {
			if element == nil {
				return
			}
			chSuccess <- element.Check(value)
		}(cronTimer.dayCondition, t.Day(), cronTimer.chSuccess)

		go func(element *ConditionElement, value int, chSuccess chan bool) {
			if element == nil {
				return
			}
			chSuccess <- element.Check(value)
		}(cronTimer.monthCondition, int(t.Month()), cronTimer.chSuccess)

		go func(element *ConditionElement, value int, chSuccess chan bool) {
			if element == nil {
				return
			}
			chSuccess <- element.Check(value)
		}(cronTimer.weekdayCondition, int(t.Weekday()), cronTimer.chSuccess)
	}(cronTimer, t)

	successCount := 0
	for {
		isSuccess := <-cronTimer.chSuccess
		if !isSuccess {
			return false
		}

		successCount++
		if successCount >= elementCount {
			return true
		}
	}
}

func (cronTimer *CronTimer) addElement(value string, maxValue int) *ConditionElement {
	data := strings.Split(value, ",")
	cElement := ConditionElement{Max: maxValue}

	for _, el := range data {
		mul := strings.Split(el, "/")
		limitComndition := strings.Split(mul[0], "-")

		min, _ := strconv.Atoi(limitComndition[0])
		max := 0
		if len(limitComndition) > 1 {
			max, _ = strconv.Atoi(limitComndition[1])
		}

		if len(mul) > 1 {
			mulValue, _ := strconv.Atoi(mul[1])
			if max == 0 {
				max = maxValue
			}

			cElement.Conditions = append(cElement.Conditions, Condition{
				Min:          min,
				Max:          max,
				Multiplicity: mulValue,
			})
			continue
		}

		if max > 0 {
			cElement.Conditions = append(cElement.Conditions, Condition{
				Min: min,
				Max: max,
			})
			continue
		}

		cElement.EnumValues = append(cElement.EnumValues, min)
	}

	return &cElement
}

// Parse parse cron record
func (cronTimer *CronTimer) Parse(value string) {
	pregMatch := `([\d\/\,\*\-]+)\s+([\d\/\,\*\-]+)\s+([\d\/\,\*\-]+)\s+([\d\/\,\*\-]+)\s+([\d\/\,\*\-]+)`

	data := regexp.MustCompile(pregMatch).FindAllStringSubmatch(value, -1)
	minute := data[0][1]
	hour := data[0][2]
	day := data[0][3]
	month := data[0][4]
	weekday := data[0][5]

	if minute != "*" {
		cronTimer.minuteCondition = cronTimer.addElement(minute, 59)
	}
	if hour != "*" {
		cronTimer.hourCondition = cronTimer.addElement(hour, 23)
	}
	if day != "*" {
		cronTimer.dayCondition = cronTimer.addElement(day, 31)
	}
	if month != "*" {
		cronTimer.monthCondition = cronTimer.addElement(month, 12)
	}
	if weekday != "*" {
		cronTimer.weekdayCondition = cronTimer.addElement(weekday, 6)
	}
}

func (cronTimer *CronTimer) getCountCondition() int {
	count := 0
	conditions := []*ConditionElement{
		cronTimer.minuteCondition,
		cronTimer.hourCondition,
		cronTimer.dayCondition,
		cronTimer.monthCondition,
		cronTimer.weekdayCondition,
	}

	for _, condition := range conditions {
		if condition != nil {
			count++
		}
	}

	return count
}

// Start start cron rule
func (cronTimer *CronTimer) Start(fnExec func(interface{}), data interface{}) {
	cronTimer.ticker = time.NewTicker(time.Minute)
	go func(fnExec func(interface{}), data interface{}) {
		for {
			if !cronTimer.Check(<-cronTimer.ticker.C) {
				continue
			}

			fnExec(data)
		}
	}(fnExec, data)
}

// NewCronTimer new cron rule
func NewCronTimer(value string) *CronTimer {
	cronTimer := &CronTimer{
		data: value,
	}
	cronTimer.Parse(value)

	return cronTimer
}
