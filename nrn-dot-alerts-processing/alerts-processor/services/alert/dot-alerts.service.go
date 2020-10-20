package dotAlertService

import (
	"bytes"
	db "dot-alerts/alerts-processor/database"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type AlertsMessage struct {
	Alerts []db.Alert
	Dot    db.Dot
}

type AlertNotificationMessage struct {
	DotID            string
	NotificationType string
	Title            string
	Body             string
}

var client = &http.Client{}

func GetDotsAlerts(sensorId string) AlertsMessage {
	var alerts AlertsMessage
	url := os.Getenv("WEB_API_HOST") + "/alert/findBySensor"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("_sensor", sensorId)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("WEB_API_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &alerts)
	if err != nil {
		panic(err)
	}

	return alerts
}

func SendAlertNotification(dotId string, alertType string, title string, body string) (*http.Response, error) {

	var notification AlertNotificationMessage

	notification.DotID = dotId
	notification.NotificationType = alertType
	notification.Title = title
	notification.Body = body

	requestBody, err := json.Marshal(notification)
	if err != nil {
		panic(err)
	}

	url := os.Getenv("WEB_API_HOST") + "/notification/sendNotification"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("WEB_API_TOKEN"))
	return client.Do(req)
}

func SaveAlertHistory(alertId string, description string) (*http.Response, error) {
	var history db.AlertHistory

	history.AlertId = alertId
	history.Description = description

	requestBody, err := json.Marshal(history)
	if err != nil {
		panic(err)
	}

	url := os.Getenv("WEB_API_HOST") + "/alertHistory"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("WEB_API_TOKEN"))
	return client.Do(req)
}

func UpdateAlert(alert db.Alert) (*http.Response, error) {

	requestBody, err := json.Marshal(alert)
	if err != nil {
		panic(err)
	}

	url := os.Getenv("WEB_API_HOST") + "/alert/" + alert.Id

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("WEB_API_TOKEN"))
	return client.Do(req)
}
