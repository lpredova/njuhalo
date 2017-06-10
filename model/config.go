package model

// Configuration struct for holding config data
type Configuration struct {
	RunIntervalMin          int     `json:"runIntervalMinutes"`
	SleepIntervalSec        int     `json:"sleepIntervalSeconds"`
	Slack                   bool    `json:"slack"`
	SlackToken              string  `json:"slackToken"`
	SlackChannelID          string  `json:"slackChannelId"`
	SlackNotificiationColor string  `json:"slackNotificationColor"`
	Mail                    bool    `json:"mail"`
	MaliTo                  string  `json:"to"`
	MailGunDomain           string  `json:"mailgunDomain"`
	MailGunAPIKey           string  `json:"mailgunAPIKey"`
	MailGunPublicKey        string  `json:"mailgunPublicKey"`
	Queries                 []Query `json:"queries,omitempty"`
}

// Query struct is used for appending new queries to config
type Query struct {
	BaseQueryPath string            `json:"baseQueryPath,omitempty"`
	Filters       map[string]string `json:"filters,omitempty"`
}
