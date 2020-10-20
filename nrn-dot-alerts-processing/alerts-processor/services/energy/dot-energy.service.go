package dotEnergyService

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type DotEnergyMessage struct {
	ActiveEnergy   []float64 `json:"active_energy"`
	ReactiveEnergy []float64 `json:"reactive_energy"`
	Timestamp      []string  `json:"timestamp"`
}

var client = &http.Client{}

func GetDotEnergy(dotId string, startTime time.Time, endTime time.Time, groupBy string) DotEnergyMessage {
	var energy DotEnergyMessage

	url := os.Getenv("WEB_API_HOST") + "/measuringPoint/getElectricEnergy"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("_point", dotId)
	q.Add("startTime", startTime.Format(time.RFC3339))
	q.Add("endTime", endTime.Format(time.RFC3339))
	q.Add("groupBy", groupBy)
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

	err = json.Unmarshal(body, &energy)
	if err != nil {
		panic(err)
	}

	return energy
}
