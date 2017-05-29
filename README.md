shNjuskhalo
===========

shNjuskhalo is simple script that watches popular online store njuskalo.hr and
notifies owner for new occurrences of wanted items.

### Usage
njuhalo s


You need to configure your watcher how often do you want to run it and for which
queries it would be active

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
      "baseQueryPath": "item/cars/path",
      "filter":{
        "maxPrice" : "11450"
        "Year" : "1992"
      }
    }
  ]
}
