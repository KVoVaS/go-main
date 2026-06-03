package domain

import "time"

type Reading struct {
	DeviceID  string
	Value     float64
	Timestamp time.Time
}
