package alert

import (
	"fmt"
	"strings"

	"github.com/lpredova/njuhalo/configuration"
	"github.com/lpredova/njuhalo/model"
	"github.com/nlopes/slack"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
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
	conf = configuration.ParseConfig()

	mailgunClient := mailgun.NewMailgun(conf.MailGunDomain, conf.MailGunAPIKey, conf.MailGunPublicKey)

	emailBody := "<ol>"
	for _, offer := range offers {
		emailBody += fmt.Sprintf("<li><img src='%s' alt='%s'><br>%s<br><a href='%s'>%s</a></li>", offer.Image, offer.Name, offer.Price, offer.URL, offer.Name)
	}
	emailBody += "<ol>"

	if len(offers) > 0 {
		message := mailgun.NewMessage(
			"njuhalo@njuh.hr",
			fmt.Sprintf("%d NEW ITEMS FOUND", len(offers)),
			emailBody,
			conf.MaliTo)

		message.SetHtml(emailBody)
		_, _, err := mailgunClient.Send(message)
		if err != nil {
			fmt.Printf(err.Error())
		} else {
			fmt.Printf("Email sent to %s", conf.MaliTo)
		}
	}
}
