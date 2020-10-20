package helpers

import (
	db "dot-alerts/alerts-processor/database"
	"encoding/json"
	"strconv"
	"time"
)

func RepeatAlert(alertId string, repeat int) bool {
	repeatCached, err := db.RClient().Get("repeat-" + alertId).Result()
	if err != nil {
		if string(err.Error()) == "redis: nil" {
			if repeat > 0 {
				cacheRepeat(alertId, repeat)
			}
			return true
		} else {
			panic(err)
		}
	} else {
		t, err := time.Parse(time.RFC3339, repeatCached)
		if err != nil {
			panic(err)
		}
		if time.Now().After(t) {
			if repeat > 0 {
				cacheRepeat(alertId, repeat)
			}
			return true
		} else {
			return false
		}
	}
}
func cacheRepeat(id string, repeat int) {
	timeCached := time.Now().Add(time.Duration(repeat) * time.Minute)
	db.RClient().Set("repeat-"+id, timeCached.Format(time.RFC3339), time.Duration(repeat+10)*time.Minute)
}

func PercentageAlert(alertId string, percent float64, duration time.Duration) bool {
	percentString := strconv.FormatFloat(percent, 'E', -1, 64)
	_, err := db.RClient().Get("percent-" + alertId + percentString).Result()
	if err != nil {
		if string(err.Error()) == "redis: nil" {
			db.RClient().Set("percent-"+alertId+percentString, "true", duration)
			return true
		} else {
			panic(err)
		}
	}
	return false
}

func WeekDayActivation(daysOfWeek string) bool {
	var alertActivationWeedDays []int
	// Alert Week Days Activation
	err := json.Unmarshal([]byte(daysOfWeek), &alertActivationWeedDays)
	if err != nil {
		panic(err)
	}
	// Week Day Today
	weekDayToday := int(time.Now().Weekday())

	// Its day of activate alert?
	_, found := Find(alertActivationWeedDays, weekDayToday)
	if found || len(alertActivationWeedDays) == 0 {
		return true
	}
	return false
}

func Find(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
