package model

// Configuration struct for holding config data
type Configuration struct {
	RunIntervalMin          int    `json:"runIntervalMinutes"`
	SleepIntervalSec        int    `json:"sleepIntervalSeconds"`
	SlackChannelID          string `json:"slackChannelId"`
	SlackNotificiationColor string `json:"slackNotificationColor"`
	Queries                 []struct {
		BaseQueryPath string            `json:"baseQueryPath,omitempty"`
		Filters       map[string]string `json:"filters,omitempty"`
	} `json:"queries,omitempty"`
}
