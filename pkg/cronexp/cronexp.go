package cronexp

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

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

// String implements standard Stringer interface
func (x CronExp) String() string {
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
	minutes, err := x.parseNumber(split[0], "minutes", MinuteRange[0], MinuteRange[1], nil)
	if err != nil {
		return err
	}
	x.minutes = minutes

	// Hours
	hours, err := x.parseNumber(split[1], "hours", HourRange[0], HourRange[1], nil)
	if err != nil {
		return err
	}
	x.hours = hours

	// Days of month
	dom, err := x.parseNumber(split[2], "days of month", DomRange[0], DomRange[1], nil)
	if err != nil {
		return err
	}
	x.daysOfMonth = dom

	// Months
	months, err := x.parseNumber(split[3], "months", MonthRange[0], MonthRange[1], MonthNames)
	if err != nil {
		return err
	}
	x.months = months

	// Days of week
	dow, err := x.parseNumber(split[4], "days of week", DowRange[0], DowRange[1], WeekNames)
	if err != nil {
		return err
	}
	x.daysOfWeek = dow

	// Command
	x.command = split[5]

	return nil
}

// parseNumber parses and validates the cron input for the given range
func (x CronExp) parseNumber(raw, fieldName string, min, max int, replaceValues map[string]int) ([]int, error) {
	// ALL
	if raw == "*" {
		return fill(min, max), nil
	}

	rawReplaced := raw
	if replaceValues != nil {
		for k, v := range replaceValues {
			rawReplaced = strings.ReplaceAll(rawReplaced, k, fmt.Sprint(v))
		}
	}

	// SIMPLE NUMBER (x)
	if simple, err := strconv.Atoi(rawReplaced); err == nil {
		if simple < min || simple > max {
			return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", fieldName, rawReplaced, min, max)
		}

		return []int{simple}, nil
	}

	// FREQUENCY (*/x) -- ignores replacements
	if sp := strings.Split(raw, "*/"); len(sp) == 2 {
		frequency, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid number", fieldName, raw)
		}

		if frequency < 1 {
			return nil, fmt.Errorf("failed to parse cron %s [%s] frequency can not be less than 1", fieldName, raw)
		}

		if frequency < min || frequency > max {
			return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", fieldName, raw, min, max)
		}

		// Assume frequency == 1 is allowed, return all
		if frequency == 1 {
			return fill(min, max), nil
		}

		// Frequency only occurs once in the set
		if frequency > max/2 {
			return []int{frequency}, nil
		}

		res := []int{min}
		counter := min
		for {
			counter += frequency
			if counter > max {
				break
			}
			res = append(res, counter)
		}

		return res, nil
	}

	// RANGE (x-y)
	if sp := strings.Split(rawReplaced, "-"); len(sp) == 2 {
		first, err := strconv.Atoi(sp[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid first number", fieldName, rawReplaced)
		}

		second, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid second number", fieldName, rawReplaced)
		}

		// Assume this is invalid
		if first >= second {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid range", fieldName, rawReplaced)
		}

		if first < min || second > max {
			return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", fieldName, rawReplaced, min, max)
		}

		return fill(first, second), nil
	}

	// LIST (x,y,z)
	if sp := strings.Split(rawReplaced, ","); len(sp) > 1 {
		deduplicator := make(map[int]bool)

		for _, item := range sp {
			num, err := strconv.Atoi(item)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cron %s [%s] invalid number", fieldName, rawReplaced)
			}

			if num < min || num > max {
				return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", fieldName, rawReplaced, min, max)
			}

			deduplicator[num] = true
		}

		res := make([]int, len(deduplicator))
		count := 0
		for k := range deduplicator {
			res[count] = k
			count++
		}

		sort.Ints(res)

		return res, nil
	}

	return nil, fmt.Errorf("failed to parse cron %s [%s] unexpected format", fieldName, raw)
}
