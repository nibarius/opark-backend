# Push notifications #

The backend can send push notifications to clients using FCM.
This document describes the format of the data sent in the push
notifications.

## Base format ##
All notifications contains a `type` parameter that tells what type of
notification it is and a `type_data` parameter with data specific to
that particular notification type.

``` json
{
  "ttl": "72000",
  "type": "available",
  "type_data": "eyJmcmVlIjoiMiJ9"
}
```

### Parameters ###

Name | Type | Mandatory | Description
---- | ---- | --------- | -----------
ttl | string | Optional | If set it specifies for how long (in seconds) the notification should be shown to the user.
type | string | Mandatory | The type of the notification
type_data | string | Mandatory | Data specific for the the given `type`. Data is base64 encoded JSON.


### Available notification types ###

Name | Description
---- | -----------
available| Notifications sent when the garage becomes available.


## available notifications ##
Notifications sent when the garage becomes available, that is goes from being
full to not being full. The notification is sent all users who have registered
on the wait list to be notified about this.

### type_data parameters ###

Name | Type | Mandatory | Description
---- | ---- | --------- | -----------
free | string | Yes | The number of free spots in the garage at the time of sending the notification

