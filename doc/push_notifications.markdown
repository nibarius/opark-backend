# Push notifications #

The backend can send push notifications to clients using FCM.
This document describes the format of the data sent in the push
notifications.

## Base format ##
All notifications contains a `type` parameter that tells what type of
notification it is. The rest of the data in the notification depends
on what type of notification it is.

``` json
{
  "ttl": "72000",
  "type": "available",
  "type_data": {
    "free": "6"
  }
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
available| Notifications sent when the garage goes from full to not full.


## available notifications ##
Notifications sent when the garage goes from full to not full. The
notification is sent all users who have registered on the wait list to
be notified about this.

### Type specific parameters ###

Name | Type | Mandatory | Description
---- | ---- | --------- | -----------
free | string | Yes | The number of free spots in the garage at the time of sending the notification

