package cronexp

import (
	"fmt"
	"strings"
)

const AllValues = "*"

/*
	https://en.wikipedia.org/wiki/Cron#CRON_expression

	Minutes 0–59
	Hours 0–23
	Day of month 1–31
	Month 1–12 or JAN–DEC
	Day of week 0–6 or SUN–SAT
*/

var (
	MinuteRange = []int{0, 59}
	HourRange   = []int{0, 23}
	DomRange    = []int{1, 31}
	MonthRange  = []int{1, 12}
	DowRange    = []int{0, 6}

	MonthNames = map[string]int{
		"JAN": 1,
		"FEB": 2,
		"MAR": 3,
		"APR": 4,
		"MAY": 5,
		"JUN": 6,
		"JUL": 7,
		"AUG": 8,
		"SEP": 9,
		"OCT": 10,
		"NOV": 11,
		"DEC": 12,
	}

	WeekNames = map[string]int{
		"SUN": 0,
		"MON": 1,
		"TUE": 2,
		"WED": 3,
		"THU": 4,
		"FRI": 5,
		"SAT": 6,
	}
)

type CronExp struct {
	minutes     []int
	hours       []int
	daysOfMonth []int
	months      []int
	daysOfWeek  []int
	command     string
}

// New creates a new CronExp object and parses the raw input string
func New(rawInput string) (*CronExp, error) {
	c := &CronExp{}
	if err := c.Parse(rawInput); err != nil {
		return nil, err
	}

	return c, nil
}

func (x CronExp) Expand() string {
	output := fmt.Sprintf(`
minute        %v
hour          %v
day of month  %v
month         %v
day of week   %v
command       %v
`, x.minutes, x.hours, x.daysOfMonth, x.months, x.daysOfWeek, x.command)

	// Trim
	output = strings.NewReplacer("[", "", "]", "").Replace(output)

	return output
}

// Parse breaks down the raw input string and parses the individual components
func (x *CronExp) Parse(rawInput string) error {
	split := strings.Split(rawInput, " ")

	if len(split) != 6 {
		return fmt.Errorf("invalid cron expression (input must contain 6 items)")
	}

	// Minutes
	minutes, err := FieldParser{
		input:            split[0],
		fieldName:        "minutes",
		minRange:         MinuteRange[0],
		maxRange:         MinuteRange[1],
		optReplaceValues: nil,
	}.GenerateValues()
	if err != nil {
		return err
	}
	x.minutes = minutes

	// Hours
	hours, err := FieldParser{
		input:            split[1],
		fieldName:        "hours",
		minRange:         HourRange[0],
		maxRange:         HourRange[1],
		optReplaceValues: nil,
	}.GenerateValues()
	if err != nil {
		return err
	}
	x.hours = hours

	// Days of month
	dom, err := FieldParser{
		input:            split[2],
		fieldName:        "days of month",
		minRange:         DomRange[0],
		maxRange:         DomRange[1],
		optReplaceValues: nil,
	}.GenerateValues()
	if err != nil {
		return err
	}
	x.daysOfMonth = dom

	// Months
	months, err := FieldParser{
		input:            split[3],
		fieldName:        "months",
		minRange:         MonthRange[0],
		maxRange:         MonthRange[1],
		optReplaceValues: MonthNames,
	}.GenerateValues()
	if err != nil {
		return err
	}
	x.months = months

	// Days of week
	dow, err := FieldParser{
		input:            split[4],
		fieldName:        "days of week",
		minRange:         DowRange[0],
		maxRange:         DowRange[1],
		optReplaceValues: WeekNames,
	}.GenerateValues()
	if err != nil {
		return err
	}
	x.daysOfWeek = dow

	// Command
	x.command = split[5]

	return nil
}
