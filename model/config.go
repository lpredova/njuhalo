package model

// Configuration struct for holding config data
type Configuration struct {
	RunIntervalMin   int     `json:"runIntervalMinutes"`
	SleepIntervalSec int     `json:"sleepIntervalSeconds"`
	Queries          []query `json:"queries"`
}

type query struct {
	BaseQueryPath string            `json:"baseQueryPath"`
	Filters       map[string]string `json:"filters"`
}
