package cronexp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFill(t *testing.T) {
	test := assert.New(t)

	{
		min := 1
		max := 12
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		actual := fill(min, max)

		test.EqualValues(expected, actual)
	}

	{
		min := 0
		max := 21
		expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}
		actual := fill(min, max)

		test.EqualValues(expected, actual)
	}

	{
		min := -2
		max := 5
		expected := []int{-2, -1, 0, 1, 2, 3, 4, 5}
		actual := fill(min, max)

		test.EqualValues(expected, actual)
	}
}

func TestCronExp_Parse_Success(t *testing.T) {
	test := assert.New(t)

	{
		input := "* * * * * /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := fill(0, 59)
				test.EqualValues(expected, cx.minutes)
			}
			{
				expected := fill(0, 23)
				test.EqualValues(expected, cx.hours)
			}
			{
				expected := fill(1, 31)
				test.EqualValues(expected, cx.daysOfMonth)
			}
			{
				expected := fill(1, 12)
				test.EqualValues(expected, cx.months)
			}
			{
				expected := fill(0, 6)
				test.EqualValues(expected, cx.daysOfWeek)
			}
		}
	}

	{
		input := "15 5 7 11 4 /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := []int{15}
				test.EqualValues(expected, cx.minutes)
			}
			{
				expected := []int{5}
				test.EqualValues(expected, cx.hours)
			}
			{
				expected := []int{7}
				test.EqualValues(expected, cx.daysOfMonth)
			}
			{
				expected := []int{11}
				test.EqualValues(expected, cx.months)
			}
			{
				expected := []int{4}
				test.EqualValues(expected, cx.daysOfWeek)
			}

		}
	}

	{
		input := "*/15 */4 */6 */3 */2 /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := []int{0, 15, 30, 45}
				test.EqualValues(expected, cx.minutes)
			}
			{
				expected := []int{0, 4, 8, 12, 16, 20}
				test.EqualValues(expected, cx.hours)
			}
			{
				expected := []int{1, 7, 13, 19, 25, 31}
				test.EqualValues(expected, cx.daysOfMonth)
			}
			{
				expected := []int{1, 4, 7, 10}
				test.EqualValues(expected, cx.months)
			}
			{
				expected := []int{0, 2, 4, 6}
				test.EqualValues(expected, cx.daysOfWeek)
			}
		}
	}

	{
		input := "3-14 11-22 9-15 5-9 3-6 /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := fill(3, 14)
				test.EqualValues(expected, cx.minutes)
			}
			{
				expected := fill(11, 22)
				test.EqualValues(expected, cx.hours)
			}
			{
				expected := fill(9, 15)
				test.EqualValues(expected, cx.daysOfMonth)
			}
			{
				expected := fill(5, 9)
				test.EqualValues(expected, cx.months)
			}
			{
				expected := fill(3, 6)
				test.EqualValues(expected, cx.daysOfWeek)
			}
		}
	}

	{
		input := "4,9,9,9,59,44,13,27,58 22,10,9,9,14,16,16,16,19 5,3,11,24,18,31,31 1,2,2,5,12,11 1,2,3 /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := []int{4, 9, 13, 27, 44, 58, 59}
				test.EqualValues(expected, cx.minutes)
			}
			{
				expected := []int{9, 10, 14, 16, 19, 22}
				test.EqualValues(expected, cx.hours)
			}
			{
				expected := []int{3, 5, 11, 18, 24, 31}
				test.EqualValues(expected, cx.daysOfMonth)
			}
			{
				expected := []int{1, 2, 5, 11, 12}
				test.EqualValues(expected, cx.months)
			}
			{
				expected := []int{1, 2, 3}
				test.EqualValues(expected, cx.daysOfWeek)
			}
		}
	}
}

func TestCronExp_Parse_WithNames(t *testing.T) {
	test := assert.New(t)

	{
		input := "* * * JUL THU /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := []int{7}
				test.EqualValues(expected, cx.months)
			}
			{
				expected := []int{4}
				test.EqualValues(expected, cx.daysOfWeek)
			}
		}
	}

	{
		input := "* * * MAY-SEP WED-SAT /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := fill(5, 9)
				test.EqualValues(expected, cx.months)
			}
			{
				expected := fill(3, 6)
				test.EqualValues(expected, cx.daysOfWeek)
			}
		}
	}

	{
		input := "* * * JAN,FEB,FEB,MAY,DEC,NOV MON,TUE,WED /usr/bin/find"
		cx, err := New(input)

		if test.Nil(err) {
			{
				expected := []int{1, 2, 5, 11, 12}
				test.EqualValues(expected, cx.months)
			}
			{
				expected := []int{1, 2, 3}
				test.EqualValues(expected, cx.daysOfWeek)
			}
		}
	}
}

func TestCronExp_ToString(t *testing.T) {
	test := assert.New(t)

	{
		input := "3-14 11-22 9-15 5-9 3-6 /usr/bin/find"
		cx, err := New(input)

		expected := `
minute        3 4 5 6 7 8 9 10 11 12 13 14
hour          11 12 13 14 15 16 17 18 19 20 21 22
day of month  9 10 11 12 13 14 15
month         5 6 7 8 9
day of week   3 4 5 6
command       /usr/bin/find
`

		if test.Nil(err) {
			actual := cx.Expand()
			test.EqualValues(expected, actual)
		}
	}
}
