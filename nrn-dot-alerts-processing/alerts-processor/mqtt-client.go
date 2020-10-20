package alertsProcessing

import (
	cover "dot-alerts/alerts-processor/circuit-overload"
	targ "dot-alerts/alerts-processor/consumption-targets"
	db "dot-alerts/alerts-processor/database"
	onOff "dot-alerts/alerts-processor/on-off"
	as "dot-alerts/alerts-processor/services/alert"
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func messageHandler(c mqtt.Client, msg mqtt.Message) {

	var onOffAlerts []db.Alert
	var targetAlerts []db.Alert

	var electricalData map[string]map[string][]interface{}

	err := json.Unmarshal(msg.Payload(), &electricalData)
	if err != nil {
		panic(err)
	}
	var alertsMessage as.AlertsMessage
	for _sensor, sensorData := range electricalData {

		dotId, err := db.RClient().Get(_sensor).Result()
		if err != nil {
			if string(err.Error()) == "redis: nil" {
				alertsMessage = as.GetDotsAlerts(_sensor)

				hasAlerts := len(alertsMessage.Alerts) != 0
				cacheInfo := alertsMessage.Dot.Id

				if hasAlerts {
					alertsData, err := json.Marshal(alertsMessage)
					if err != nil {
						panic(err)
					}
					db.RClient().Set(alertsMessage.Dot.Id, string(alertsData), 7*24*time.Hour)
				} else {
					cacheInfo = "has-no-alerts"
				}

				db.RClient().Set(_sensor, cacheInfo, 7*24*time.Hour)

				for i := 0; i < len(alertsMessage.Dot.Sensors); i++ {
					db.RClient().Set(alertsMessage.Dot.Sensors[i].Id, cacheInfo, 7*24*time.Hour)
				}
			} else {
				panic(err)
			}
		}
		if dotId == "has-no-alerts" {
			break
		}
		if dotId != "" {
			a, err := db.RClient().Get(dotId).Result()
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal([]byte(a), &alertsMessage)
			if err != nil {
				panic(err)
			}
		}

		var now = time.Now()
		for _, alert := range alertsMessage.Alerts {
			onlyTime := time.Date(alert.NotifTimeBegin.Year(), alert.NotifTimeBegin.Month(), alert.NotifTimeBegin.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)
			if (alert.NotifTimeBegin == time.Time{}) || onlyTime.After(alert.NotifTimeBegin) && onlyTime.Before(alert.NotifTimeEnd) {
				if alert.AlertType == "on-off" {
					onOffAlerts = append(onOffAlerts, alert)
				} else if alert.AlertType == "consumption-target" {
					targetAlerts = append(targetAlerts, alert)
				}
			}
		}

		// Process OnOFf
		if len(onOffAlerts) != 0 {
			go onOff.OnOffProcessing(alertsMessage.Dot, onOffAlerts, sensorData)
		}
		// Process Consumption Target
		if len(targetAlerts) != 0 {
			go targ.ConsumptionTargetProcessing(alertsMessage.Dot, targetAlerts, sensorData)
		}
		// Process Circuit Overload
		go cover.CircuitOverloadProcessing(alertsMessage.Dot, sensorData)
	}
}

func connLostHandler(c mqtt.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)

	//Perform additional action...
}

func MQTTClient() {
	//create a ClientOptions
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://" + os.Getenv("MQTT_BROKER_HOST") + ":" + os.Getenv("MQTT_BROKER_PORT")).
		SetClientID("nrn-onoff-alerts").
		SetUsername(os.Getenv("MQTT_BROKER_USER")).
		SetPassword(os.Getenv("MQTT_BROKER_PASSWORD")).
		SetDefaultPublishHandler(messageHandler).
		SetConnectionLostHandler(connLostHandler)

	//set OnConnect handler as anonymous function
	//after connected, subscribe to topic
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Printf("Client Dot Alerts connected\n")

		//Subscribe here, otherwise after connection lost,
		//you may not receive any message
		if token := c.Subscribe(os.Getenv("MQTT_ELECTRICAL_DATA_TOPIC"), 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}

	//create and start a client using the above ClientOptions
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		//Lazy...
		time.Sleep(500 * time.Millisecond)
	}
}
