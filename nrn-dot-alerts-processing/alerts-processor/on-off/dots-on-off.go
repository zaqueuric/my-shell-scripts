package onOff

import (
	db "dot-alerts/alerts-processor/database"
	ah "dot-alerts/alerts-processor/helpers"
	as "dot-alerts/alerts-processor/services/alert"
	"strconv"

	stat "gonum.org/v1/gonum/stat"
)

func OnOffProcessing(dot db.Dot, alerts []db.Alert, sensorData map[string][]interface{}) {

	var dotState string = "offline"

	// Power slice => interface to float64
	power := make([]float64, len(sensorData["power"]))
	for i, p := range sensorData["power"] {
		power[i] = p.(float64)
	}

	// Get Measured Power Mean
	meanPower := stat.Mean(power, nil)

	// Get Power State Threshold
	powerStateThreshold, err := strconv.ParseFloat(dot.TechnicalInfo["powerStateThreshold"], 32)
	if err != nil {
		powerStateThreshold = 60.0
	}

	// Current Dot State (ON or OFF)
	if meanPower <= powerStateThreshold {
		dotState = "off"
	} else {
		dotState = "on"
	}

	// Alerts processing
	for _, alert := range alerts {
		if alert.IsEnabled && ah.RepeatAlert(alert.Id, alert.Repeat) && ah.WeekDayActivation(alert.DaysOfWeek) {
			// State for activation?
			if dotState == alert.Attributes["activationState"] {
				_, err := as.SendAlertNotification(dot.Id, alert.AlertType, dot.Name, alert.Description)
				if err != nil {
					panic(err)
				}
				_, err = as.SaveAlertHistory(alert.Id, alert.Description)
				if err != nil {
					panic(err)
				}
				if alert.DaysOfWeek == "[]" {
					alert.IsEnabled = false
					_, err = as.UpdateAlert(alert)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
