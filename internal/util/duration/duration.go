package duration

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var timeNow = time.Now

func Parse(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty duration string")
	}

	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}

	now := timeNow()

	if strings.HasSuffix(s, "mo") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "mo"))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration: %s", s)
		}
		return now.AddDate(0, -n, 0), nil
	}

	if strings.HasSuffix(s, "y") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "y"))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration: %s", s)
		}
		return now.AddDate(-n, 0, 0), nil
	}

	if strings.HasSuffix(s, "w") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "w"))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration: %s", s)
		}
		return now.AddDate(0, 0, -n*7), nil
	}

	if strings.HasSuffix(s, "d") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration: %s", s)
		}
		return now.AddDate(0, 0, -n), nil
	}

	return time.Time{}, fmt.Errorf("invalid duration: %s", s)
}
