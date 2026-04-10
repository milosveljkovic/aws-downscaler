package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func daysInRange(startDay, endDay int) []int {
	if startDay > 6 || startDay < 0 || endDay > 6 || endDay < 0 {
		return nil
	}
	d := startDay
	days := []int{}
	for d != endDay {
		d = (d + 1) % 7
		days = append(days, d)
	}
	return days
}

func ParseToSeconds(input string) int64 {
	defaultSeconds := int64(60)

	// Try standard duration parsing first
	d, err := time.ParseDuration(input)
	if err == nil {
		s := int64(d.Seconds())
		if s < defaultSeconds {
			log.Warn().Msgf("Interval need to be more than %d", defaultSeconds)
			log.Info().Msgf("Setting default interval of %ds", defaultSeconds)
			return defaultSeconds
		}
		return s
	}

	// If parsing failed, try plain number (assume seconds)
	if n, err := strconv.Atoi(input); err == nil {
		s := int64(n)
		if s < defaultSeconds {
			log.Debug().Msgf("Interval need to be more than %d", defaultSeconds)
			log.Info().Msgf("Setting default interval of %ds", defaultSeconds)
			return defaultSeconds
		}
		return s
	}

	// Fallback to default
	return defaultSeconds
}

var dayToIndex = map[string]int{
	"Mon": 1,
	"Tue": 2,
	"Wed": 3,
	"Thu": 4,
	"Fri": 5,
	"Sat": 6,
	"Sun": 7,
}

type rule struct {
	startDay    int
	endDay      int
	startMinute int
	endMinute   int
}

func IsNowInDowntime(downtime string) (bool, error) {
	return isInDowntimeAt(downtime, time.Now())
}

func isInDowntimeAt(downtime string, now time.Time) (bool, error) {
	rules, loc, err := parseDowntime(downtime)
	if err != nil {
		return false, err
	}

	localNow := now.In(loc)
	log.Info().Msgf("Localtime: %s %s", localNow.Format("2006-01-02 15:04:05"), loc.String())
	log.Info().Msgf("Downtime: %s", downtime)
	currentDay := goWeekdayToIndex(localNow.Weekday())
	currentMinute := localNow.Hour()*60 + localNow.Minute()

	for _, r := range rules {
		if matchesRule(r, currentDay, currentMinute) {
			log.Info().Msg("downtime: TRUE")
			return true, nil
		}
	}

	log.Info().Msg("downtime: FALSE")
	return false, nil
}

func parseDowntime(downtime string) ([]rule, *time.Location, error) {
	parts := strings.Split(downtime, ",")
	if len(parts) == 0 {
		return nil, nil, fmt.Errorf("empty downtime spec")
	}

	var rules []rule
	var timezoneName string
	var loc *time.Location

	for _, raw := range parts {
		part := strings.TrimSpace(raw)
		if part == "" {
			continue
		}

		fields := strings.Fields(part)
		if len(fields) != 3 {
			return nil, nil, fmt.Errorf("invalid downtime part %q, expected: <days-range> <time-range> <timezone>", part)
		}

		dayExpr := fields[0]
		timeExpr := fields[1]
		tzName := fields[2]

		if timezoneName == "" {
			timezoneName = tzName

			var err error
			loc, err = time.LoadLocation(tzName)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid timezone %q: %w", tzName, err)
			}
		} else if tzName != timezoneName {
			return nil, nil, fmt.Errorf("mixed timezones are not allowed: found %q and %q", timezoneName, tzName)
		}

		startDay, endDay, err := parseDayRange(dayExpr)
		if err != nil {
			return nil, nil, err
		}

		startMinute, endMinute, err := parseTimeRange(timeExpr)
		if err != nil {
			return nil, nil, err
		}

		rules = append(rules, rule{
			startDay:    startDay,
			endDay:      endDay,
			startMinute: startMinute,
			endMinute:   endMinute,
		})
	}

	if len(rules) == 0 {
		return nil, nil, fmt.Errorf("empty downtime spec")
	}

	return rules, loc, nil
}

func matchesRule(r rule, currentDay, currentMinute int) bool {
	days := expandDayRange(r.startDay, r.endDay)

	// Full day, e.g. 00:00-24:00
	if r.startMinute == 0 && r.endMinute == 1440 {
		return containsDay(days, currentDay)
	}

	// Same-day window, e.g. 09:00-17:00
	if r.startMinute < r.endMinute {
		return containsDay(days, currentDay) &&
			currentMinute >= r.startMinute &&
			currentMinute < r.endMinute
	}

	// Overnight window, e.g. Mon-Fri 18:00-08:00
	// Means:
	// Mon 18:00 -> Tue 08:00
	// Tue 18:00 -> Wed 08:00
	// Wed 18:00 -> Thu 08:00
	// Thu 18:00 -> Fri 08:00
	if len(days) == 1 {
		if currentDay == days[0] && currentMinute >= r.startMinute {
			return true
		}
		if currentDay == nextDay(days[0]) && currentMinute < r.endMinute {
			return true
		}
		return false
	}

	for i := 0; i < len(days)-1; i++ {
		startD := days[i]
		endD := days[i+1]

		if currentDay == startD && currentMinute >= r.startMinute {
			return true
		}
		if currentDay == endD && currentMinute < r.endMinute {
			return true
		}
	}

	return false
}

func parseDayRange(s string) (int, int, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid day range %q", s)
	}

	start, ok := dayToIndex[parts[0]]
	if !ok {
		return 0, 0, fmt.Errorf("invalid start day %q", parts[0])
	}

	end, ok := dayToIndex[parts[1]]
	if !ok {
		return 0, 0, fmt.Errorf("invalid end day %q", parts[1])
	}

	return start, end, nil
}

func parseTimeRange(s string) (int, int, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time range %q", s)
	}

	start, err := parseHHMM(parts[0])
	if err != nil {
		return 0, 0, err
	}

	end, err := parseHHMM(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

func parseHHMM(s string) (int, error) {
	if s == "24:00" {
		return 1440, nil
	}

	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, fmt.Errorf("invalid time %q", s)
	}

	return t.Hour()*60 + t.Minute(), nil
}

func expandDayRange(start, end int) []int {
	var days []int
	d := start
	for {
		days = append(days, d)
		if d == end {
			break
		}
		d = nextDay(d)
	}
	return days
}

func nextDay(day int) int {
	if day == 7 {
		return 1
	}
	return day + 1
}

func containsDay(days []int, day int) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}

func goWeekdayToIndex(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}
