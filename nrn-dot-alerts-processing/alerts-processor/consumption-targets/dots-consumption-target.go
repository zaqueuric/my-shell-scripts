package consumptionTarget

import (
	db "dot-alerts/alerts-processor/database"
	ah "dot-alerts/alerts-processor/helpers"
	as "dot-alerts/alerts-processor/services/alert"
	es "dot-alerts/alerts-processor/services/energy"
	"encoding/json"
	"strconv"
	"time"
)

func ConsumptionTargetProcessing(dot db.Dot, alerts []db.Alert, sensorData map[string][]interface{}) {
	var activationPercentages []float64
	var percentDuration time.Duration

	// Alerts processing
	for _, alert := range alerts {
		if alert.IsEnabled && ah.WeekDayActivation(alert.DaysOfWeek) {
			endTime := time.Now()
			startTime := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, time.UTC)

			switch frequency := alert.Attributes["frequency"]; frequency {
			case "daily":
				percentDuration = time.Duration((60 - endTime.Minute())) * time.Minute
			case "byTurn":
				startTurn, err := time.Parse(time.RFC3339, alert.Attributes["startTurn"])
				if err != nil {
					panic(err)
				}
				startTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), startTurn.Hour(), startTurn.Minute(), startTurn.Second(), startTurn.Nanosecond(), time.UTC)
				endTurn, err := time.Parse(time.RFC3339, alert.Attributes["endTurn"])
				if err != nil {
					panic(err)
				}
				endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), endTurn.Hour(), endTurn.Minute(), endTurn.Second(), endTurn.Nanosecond(), time.UTC)
				percentDuration = time.Duration(endTime.Minute()) * time.Minute
			case "weekly":
				weekDay := int(endTime.Weekday())
				startTime = startTime.AddDate(0, 0, -weekDay)
				percentDuration = time.Duration((7-weekDay)*24) * time.Hour
			default:
				startTime = startTime.AddDate(0, 0, -endTime.Day()+1)
				y, m, d := endTime.Date()
				lastday := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)
				daysEndMon := lastday.Day() - d
				percentDuration = time.Duration(daysEndMon*24) * time.Hour
			}
			// Get consumption
			energy := es.GetDotEnergy(dot.Id, startTime, endTime, "total")

			target, err := strconv.ParseFloat(alert.Attributes["target"], 32)
			if err != nil {
				panic(err)
			}
			energyPercentage := energy.ActiveEnergy[0] / target

			err = json.Unmarshal([]byte(alert.Attributes["percentages"]), &activationPercentages)
			if err != nil {
				panic(err)
			}

			if len(activationPercentages) == 1 && energyPercentage > 1 && ah.PercentageAlert(alert.Id, 1, percentDuration) {
				// Send notification
				as.SendAlertNotification(dot.Id, alert.AlertType, dot.Name, alert.Description)
			} else {
				for _, p := range activationPercentages {
					if energyPercentage > p && ah.PercentageAlert(alert.Id, p, percentDuration) {
						//Send notification
						as.SendAlertNotification(dot.Id, alert.AlertType, dot.Name, alert.Description)
					}
				}
			}
		}
	}
}
