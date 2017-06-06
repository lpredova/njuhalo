package model

// Configuration struct for holding config data
type Configuration struct {
	RunIntervalMin          int    `json:"runIntervalMinutes"`
	SleepIntervalSec        int    `json:"sleepIntervalSeconds"`
	Slack                   bool   `json:"slack"`
	SlackToken              string `json:"slackToken"`
	SlackChannelID          string `json:"slackChannelId"`
	SlackNotificiationColor string `json:"slackNotificationColor"`
	Mail                    bool   `json:"mail"`
	MaliTo                  string `json:"to"`
	MailGunDomain           string `json:"mailgunDomain"`
	MailGunAPIKey           string `json:"mailgunAPIKey"`
	MailGunPublicKey        string `json:"mailgunPublicKey"`
	Queries                 []struct {
		BaseQueryPath string            `json:"baseQueryPath,omitempty"`
		Filters       map[string]string `json:"filters,omitempty"`
	} `json:"queries,omitempty"`
}
