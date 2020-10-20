package circuitOverload

import (
	db "dot-alerts/alerts-processor/database"
	as "dot-alerts/alerts-processor/services/alert"
	"strconv"

	floats "gonum.org/v1/gonum/floats"
)

func CircuitOverloadProcessing(dot db.Dot, sensorData map[string][]interface{}) {
	// Current slice => interface to float64
	current := make([]float64, len(sensorData["current"]))
	for i, p := range sensorData["current"] {
		current[i] = p.(float64)
	}

	maxCurrent := floats.Max(current)
	currentString := strconv.FormatFloat(maxCurrent, 'f', 2, 64)

	if _, ok := dot.TechnicalInfo["circuitBreakValue"]; ok {
		circuitBreak, err := strconv.ParseFloat(dot.TechnicalInfo["circuitBreakValue"], 32)
		if err != nil {
			panic(err)
		}
		// State for activation?
		if maxCurrent > 0.9*circuitBreak {
			alertType := "overload-circuit"
			title := "Sobrecarga no circuito"
			body := "A corrente no circuito de " + dot.Name + " está com valor de " + currentString + "A. O valor do disjuntor para esse ponto supervisionado é de " + dot.TechnicalInfo["circuitBreakValue"] + "A."
			as.SendAlertNotification(dot.Id, alertType, title, body)
		}
	}
}
