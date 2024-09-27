package config

import (
	"log"
	"time"
)

func SetTimeZone(timeZone string) {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	time.Local = location
	log.Printf("Timezone set to: %s", timeZone)
}
func FormatTime(t time.Time) string {
	return t.Format("2006/01/02 15:04:05")
}
