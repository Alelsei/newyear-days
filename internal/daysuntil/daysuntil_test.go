package daysuntil

import (
	"testing"
	"time"
)

func TestDaysUntilNewYear(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   time.Time
		want int
	}{
		{
			name: "start of a leap year",
			in:   time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: 366, // 2024 is a leap year, so a full 2024 remains.
		},
		{
			name: "start of a non-leap year",
			in:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: 365,
		},
		{
			name: "end of a leap year",
			in:   time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
			want: 1,
		},
		{
			name: "end of a non-leap year",
			in:   time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
			want: 1,
		},
		{
			name: "end of year, with a non-midnight time of day",
			in:   time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC),
			want: 1, // time-of-day must be ignored, only the date matters.
		},
		{
			name: "day before Feb 29 in a leap year",
			in:   time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC),
			want: 308,
		},
		{
			name: "Feb 29 itself, in a leap year",
			in:   time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			want: 307,
		},
		{
			name: "day after Feb 29 in a leap year",
			in:   time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
			want: 306,
		},
		{
			name: "same calendar day (Feb 28) in a non-leap year",
			in:   time.Date(2023, time.February, 28, 0, 0, 0, 0, time.UTC),
			want: 307,
		},
		{
			name: "day after Feb 28 in a non-leap year (no Feb 29 exists)",
			in:   time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
			want: 306,
		},
		{
			name: "an ordinary mid-year date",
			in:   time.Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC),
			want: 200,
		},
		{
			name: "input given in a non-UTC location is still handled correctly",
			in:   time.Date(2024, time.December, 31, 12, 0, 0, 0, time.FixedZone("UTC+3", 3*60*60)),
			want: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DaysUntilNewYear(tt.in)
			if got != tt.want {
				t.Errorf("DaysUntilNewYear(%s) = %d, want %d", tt.in.Format(time.RFC3339), got, tt.want)
			}
		})
	}
}

func TestDaysUntilNewYear_AlwaysPositive(t *testing.T) {
	t.Parallel()

	dates := []time.Time{
		time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2100, time.December, 31, 0, 0, 0, 0, time.UTC),
		time.Date(1999, time.July, 4, 0, 0, 0, 0, time.UTC),
	}

	for _, d := range dates {
		if got := DaysUntilNewYear(d); got <= 0 {
			t.Errorf("DaysUntilNewYear(%s) = %d, want a positive number of days", d.Format(time.RFC3339), got)
		}
	}
}
