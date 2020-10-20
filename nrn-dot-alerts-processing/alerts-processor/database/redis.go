package database

import (
	"os"
	"time"

	"github.com/go-redis/redis"
)

type Sensor struct {
	Id string `json:"_id"`
}

type Dot struct {
	Id            string            `json:"_id"`
	Name          string            `json:"name"`
	TechnicalInfo map[string]string `json:"technicalInfo"`
	Sensors       []Sensor          `json:"sensors"`
}
type Alert struct {
	Id             string            `json:"_id"`
	Description    string            `json:"description"`
	AlertType      string            `json:"type"`
	NotifTimeBegin time.Time         `json:"notifTimeBegin"`
	NotifTimeEnd   time.Time         `json:"notifTimeEnd"`
	Repeat         int               `json:"repeat"`
	DaysOfWeek     string            `json:"daysOfWeek"`
	IsEnabled      bool              `json:"isEnabled"`
	Attributes     map[string]string `json:"attributes"`
}

type AlertHistory struct {
	Description string    `json:"description"`
	AlertId     string    `json:"_alert"`
	Timestamp   time.Time `json:"timestamp"`
}

func RClient() *redis.Client {
	url := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password, // no password set
		DB:       0,
	})
	return client
}
