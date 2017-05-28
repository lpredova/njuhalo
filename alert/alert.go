package alert

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env"
	"github.com/lpredova/shnjuskhalo/configuration"
	"github.com/lpredova/shnjuskhalo/model"
	"github.com/nlopes/slack"
)

type slackConfig struct {
	SlackToken string `env:"SHNJUSKHALO_SLACK_TOKEN" envDefault:"myPreciousToken"`
}

var slackConf slackConfig
var conf model.Configuration

// SendItemsToSlack method that formats and send messages to slack
func SendItemsToSlack(offers []model.Offer) {
	env.Parse(&slackConf)

	api := slack.New(slackConf.SlackToken)
	params := slack.PostMessageParameters{}
	conf = configuration.ParseConfig()

	var attachments []slack.Attachment
	for _, offer := range offers {
		attachments = append(attachments, slack.Attachment{
			Title:     fmt.Sprintf("%s:%s", offer.ID, offer.Name),
			TitleLink: offer.URL,
			ImageURL:  offer.Image,
			Text:      fmt.Sprintf("%s %s", offer.Price, strings.TrimSpace(offer.Description)),
			Color:     conf.SlackNotificiationColor,
		})
	}
	params.Attachments = attachments
	channelID, timestamp, err := api.PostMessage(conf.SlackChannelID, "", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func sentItemsToMail() {

}
