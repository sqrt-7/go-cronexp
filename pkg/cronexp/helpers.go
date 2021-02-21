package cronexp

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// fill fills an []int slice within the given range
func fill(min, max int) []int {
	res := make([]int, max-min+1)
	for i := 0; i <= max-min; i++ {
		res[i] = i + min
	}

	return res
}

type FieldParser struct {
	input            string         // raw input string
	fieldName        string         // field name to use in error messages
	minRange         int            // min value of allowed range for field
	maxRange         int            // max value of allowed range for field
	optReplaceValues map[string]int // strings to replace if the field allows string aliases
}

// GenerateValues parses and validates the cron input segment
// and generates the array of values
func (parser FieldParser) GenerateValues() ([]int, error) {
	// ALL (*)
	if parser.input == AllValues {
		return fill(parser.minRange, parser.maxRange), nil
	}

	inputReplaced := parser.input
	if parser.optReplaceValues != nil {
		for k, v := range parser.optReplaceValues {
			inputReplaced = strings.ReplaceAll(inputReplaced, k, fmt.Sprint(v))
		}
	}

	// SIMPLE NUMBER (x)
	if simple, err := strconv.Atoi(inputReplaced); err == nil {
		if simple < parser.minRange || simple > parser.maxRange {
			return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", parser.fieldName, inputReplaced, parser.minRange, parser.maxRange)
		}

		return []int{simple}, nil
	}

	// FREQUENCY (*/x) -- ignores replacements
	if sp := strings.Split(parser.input, "*/"); len(sp) == 2 {
		frequency, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid number", parser.fieldName, parser.input)
		}

		if frequency < 1 {
			return nil, fmt.Errorf("failed to parse cron %s [%s] frequency can not be less than 1", parser.fieldName, parser.input)
		}

		if frequency < parser.minRange || frequency > parser.maxRange {
			return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", parser.fieldName, parser.input, parser.minRange, parser.maxRange)
		}

		// Assume frequency == 1 is allowed, return all
		if frequency == 1 {
			return fill(parser.minRange, parser.maxRange), nil
		}

		// Frequency only occurs once in the set
		if frequency > parser.maxRange/2 {
			return []int{frequency}, nil
		}

		res := []int{parser.minRange}
		counter := parser.minRange
		for {
			counter += frequency
			if counter > parser.maxRange {
				break
			}
			res = append(res, counter)
		}

		return res, nil
	}

	// RANGE (x-y)
	if sp := strings.Split(inputReplaced, "-"); len(sp) == 2 {
		first, err := strconv.Atoi(sp[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid first number", parser.fieldName, inputReplaced)
		}

		second, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid second number", parser.fieldName, inputReplaced)
		}

		// Assume this is invalid
		if first >= second {
			return nil, fmt.Errorf("failed to parse cron %s [%s] invalid range", parser.fieldName, inputReplaced)
		}

		if first < parser.minRange || second > parser.maxRange {
			return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", parser.fieldName, inputReplaced, parser.minRange, parser.maxRange)
		}

		return fill(first, second), nil
	}

	// LIST (x,y,z)
	if sp := strings.Split(inputReplaced, ","); len(sp) > 1 {
		deduplicator := make(map[int]bool)

		for _, item := range sp {
			num, err := strconv.Atoi(item)
			if err != nil {
				return nil, fmt.Errorf("failed to parse cron %s [%s] invalid number", parser.fieldName, inputReplaced)
			}

			if num < parser.minRange || num > parser.maxRange {
				return nil, fmt.Errorf("failed to parse cron %s [%s] (min: %d max: %d)", parser.fieldName, inputReplaced, parser.minRange, parser.maxRange)
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

	return nil, fmt.Errorf("failed to parse cron %s [%s] unexpected format", parser.fieldName, parser.input)
}
