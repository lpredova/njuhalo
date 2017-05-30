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
  "slack": false,
  "slackChannelId": "",
  "slackNotificationColor": "",
  "mail": false,
  "to": ""
  "queries" : [
    {
      "baseQueryPath": "item/rooms/path",
      "filter":{
        "maxPrice" : "5000"
        "Floors" : "4"
      }
    },
    {
      "baseQueryPath": "iznajmljivanje-stanova/zagreb",
      "filters":{
        "locationId": "2619",
        "price[max]": "260",
        "mainArea[max]": "50"
      }
    }
  ]
}


### Build

go build -ldflags -s
