package alert

import (
	"fmt"
	"strings"

	"github.com/lpredova/shnjuskhalo/configuration"
	"github.com/lpredova/shnjuskhalo/model"
	"github.com/nlopes/slack"
)

var conf model.Configuration

// SendItemsToSlack method that formats and send messages to slack
func SendItemsToSlack(offers []model.Offer) {
	conf = configuration.ParseConfig()

	api := slack.New(conf.SlackToken)
	params := slack.PostMessageParameters{}

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

// SendItemsToMail method that formats and sends mail to user
func SendItemsToMail(offers []model.Offer) {
}
