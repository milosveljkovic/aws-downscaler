package util

import (
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestDaysInRange(t *testing.T) {

	type TestDaysInRange struct {
		expLength int
		expDays   []int
		startDay  int
		endDay    int
	}

	tbd := []TestDaysInRange{
		{
			expLength: 3,
			expDays:   []int{1, 2, 3},
			startDay:  0,
			endDay:    3,
		},
		{
			expLength: 4,
			expDays:   []int{3, 4, 5, 6},
			startDay:  2,
			endDay:    6,
		},
		{
			expLength: 4,
			expDays:   []int{0, 1, 2, 3},
			startDay:  6,
			endDay:    3,
		},
	}

	for _, test := range tbd {
		days := daysInRange(test.startDay, test.endDay)
		if len(days) != test.expLength {
			t.Errorf("Expected number of days %d, got %d", test.expLength, len(days))
		}
		for _, day := range days {
			if !slices.Contains(test.expDays, day) {
				t.Errorf("Expected days '%d', got days '%d'. Seems like %d day is missing", test.expDays, days, day)
			}
		}
	}
}

func TestNextDay(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{
			name: "monday to tuesday",
			in:   1,
			want: 2,
		},
		{
			name: "friday to saturday",
			in:   5,
			want: 6,
		},
		{
			name: "saturday to sunday",
			in:   6,
			want: 7,
		},
		{
			name: "sunday wraps to monday",
			in:   7,
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nextDay(tt.in)
			if got != tt.want {
				t.Fatalf("nextDay(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestExpandDayRange(t *testing.T) {
	tests := []struct {
		name  string
		start int
		end   int
		want  []int
	}{
		{
			name:  "single day",
			start: 3,
			end:   3,
			want:  []int{3},
		},
		{
			name:  "simple forward range",
			start: 1,
			end:   5,
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "full week",
			start: 1,
			end:   7,
			want:  []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:  "wrap around friday to tuesday",
			start: 5,
			end:   2,
			want:  []int{5, 6, 7, 1, 2},
		},
		{
			name:  "wrap around sunday to tuesday",
			start: 7,
			end:   2,
			want:  []int{7, 1, 2},
		},
		{
			name:  "wrap around saturday to monday",
			start: 6,
			end:   1,
			want:  []int{6, 7, 1},
		},
		{
			name:  "reverse almost full week",
			start: 7,
			end:   6,
			want:  []int{7, 1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandDayRange(tt.start, tt.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expandDayRange(%d, %d) = %v, want %v", tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestIsInDowntimeAt(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		now     time.Time
		want    bool
		wantErr bool
	}{
		{
			name: "weekday overnight in downtime evening UTC",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 2, 19, 0, 0, 0, time.UTC), // Monday 19:00 UTC
			want: true,
		},
		{
			name: "weekday overnight in downtime next morning UTC",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 3, 7, 59, 0, 0, time.UTC), // Tuesday 07:59 UTC
			want: true,
		},
		{
			name: "weekday overnight outside downtime during day UTC",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC), // Tuesday 12:00 UTC
			want: false,
		},
		{
			name: "weekday overnight starts exactly at boundary",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 4, 18, 0, 0, 0, time.UTC), // Wednesday 18:00 UTC
			want: true,
		},
		{
			name: "weekday overnight ends exactly at boundary exclusive",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 5, 8, 0, 0, 0, time.UTC), // Thursday 08:00 UTC
			want: false,
		},
		{
			name: "saturday full day is downtime",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 7, 15, 30, 0, 0, time.UTC), // Saturday 15:30 UTC
			want: true,
		},
		{
			name: "sunday late night still downtime",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 8, 23, 59, 0, 0, time.UTC), // Sunday 23:59 UTC
			want: true,
		},
		{
			name: "monday before weekday downtime starts is not downtime",
			spec: "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 UTC",
			now:  time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC), // Monday 10:00 UTC
			want: false,
		},
		{
			name: "cross week range friday to tuesday thursday night",
			spec: "Fri-Tue 20:00-05:00 UTC",
			now:  time.Date(2026, 3, 5, 22, 0, 0, 0, time.UTC), // Thursday 22:00 UTC
			want: false,
		},
		{
			name: "cross week range friday to tuesday saturday night",
			spec: "Fri-Tue 20:00-05:00 UTC",
			now:  time.Date(2026, 3, 7, 22, 0, 0, 0, time.UTC), // Saturday 22:00 UTC
			want: true,
		},
		{
			name: "cross week range friday to tuesday monday early morning",
			spec: "Fri-Tue 20:00-05:00 UTC",
			now:  time.Date(2026, 3, 9, 4, 30, 0, 0, time.UTC), // Monday 04:30 UTC
			want: true,
		},
		{
			name: "cross week range friday to tuesday tuesday evening excluded",
			spec: "Fri-Tue 20:00-05:00 UTC",
			now:  time.Date(2026, 3, 10, 22, 0, 0, 0, time.UTC), // Tuesday 22:00 UTC
			want: false,
		},
		{
			name: "single day overnight same day evening",
			spec: "Mon-Mon 20:00-05:00 UTC",
			now:  time.Date(2026, 3, 2, 21, 0, 0, 0, time.UTC), // Monday 21:00 UTC
			want: true,
		},
		{
			name: "single day overnight next day early morning",
			spec: "Mon-Mon 20:00-05:00 UTC",
			now:  time.Date(2026, 3, 3, 4, 59, 0, 0, time.UTC), // Tuesday 04:59 UTC
			want: true,
		},
		{
			name: "single day overnight next day after end",
			spec: "Mon-Mon 20:00-05:00 UTC",
			now:  time.Date(2026, 3, 3, 5, 0, 0, 0, time.UTC), // Tuesday 05:00 UTC
			want: false,
		},
		{
			name: "same day non overnight inside window",
			spec: "Mon-Fri 09:00-17:00 UTC",
			now:  time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC), // Wednesday 10:00 UTC
			want: true,
		},
		{
			name: "same day non overnight outside window",
			spec: "Mon-Fri 09:00-17:00 UTC",
			now:  time.Date(2026, 3, 4, 18, 0, 0, 0, time.UTC), // Wednesday 18:00 UTC
			want: false,
		},
		{
			name: "timezone respected America New York inside downtime",
			spec: "Mon-Fri 18:00-08:00 America/New_York",
			now:  time.Date(2026, 3, 3, 1, 0, 0, 0, time.UTC), // Monday 20:00 in New York (EST)
			want: true,
		},
		{
			name: "timezone respected America New York outside downtime",
			spec: "Mon-Fri 18:00-08:00 America/New_York",
			now:  time.Date(2026, 3, 3, 18, 0, 0, 0, time.UTC), // Tuesday 13:00 in New York (EST)
			want: false,
		},
		{
			name:    "mixed timezones returns error",
			spec:    "Mon-Fri 18:00-08:00 UTC, Sat-Sun 00:00-24:00 America/New_York",
			now:     time.Date(2026, 3, 7, 12, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name:    "invalid timezone returns error",
			spec:    "Mon-Fri 18:00-08:00 Not_A_Real_TZ",
			now:     time.Date(2026, 3, 2, 19, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name:    "invalid day returns error",
			spec:    "Foo-Fri 18:00-08:00 UTC",
			now:     time.Date(2026, 3, 2, 19, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name:    "invalid time returns error",
			spec:    "Mon-Fri 25:00-08:00 UTC",
			now:     time.Date(2026, 3, 2, 19, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name:    "invalid format returns error",
			spec:    "Mon-Fri 18:00-08:00",
			now:     time.Date(2026, 3, 2, 19, 0, 0, 0, time.UTC),
			wantErr: true,
		},
		{
			name:    "empty spec returns error",
			spec:    "",
			now:     time.Date(2026, 3, 2, 19, 0, 0, 0, time.UTC),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isInDowntimeAt(tt.spec, tt.now)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("IsInDowntimeAt(%q, %v) error = nil, want error", tt.spec, tt.now)
				}
				return
			}

			if err != nil {
				t.Fatalf("IsInDowntimeAt(%q, %v) unexpected error: %v", tt.spec, tt.now, err)
			}

			if got != tt.want {
				t.Fatalf("IsInDowntimeAt(%q, %v) = %v, want %v", tt.spec, tt.now, got, tt.want)
			}
		})
	}
}
