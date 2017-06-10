new(/njuː/)halo
===========

newhalo is simple script written in Golang. It watches popular online store njuskalo.hr and
notifies owner for new occurrences of wanted items.

### Usage

You need to configure your watcher how often do you want to run it and for which
queries it would be active.

{
	"runIntervalMinutes": 1,
	"sleepIntervalSeconds": 2,
	"queries": [{
		"baseQueryPath": "iznajmljivanje-stanova/zagreb",
		"filters": {
			"locationId": "2619",
			"price[max]": "260",
			"mainArea[max]": "50"
		}
	}],
	"slack" : true,
	"slackToken": "",
	"slackChannelId" : "",
	"slackNotificationColor": "#fdcd00",
	"mail": true,
	"to": "myemail@gmail.com",
	"mailgunDomain":"",
	"mailgunAPIKey":"",
	"mailgunPublicKey":""
}



### Build

go build -ldflags -s
