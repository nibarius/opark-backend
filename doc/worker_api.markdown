# Worker API v1 #

## Endpoint ##
Root endpoint: https://opera-park.appspot.com/worker/v1/

## Authentication ##

No authentication is required, but only the appengine can call the worker api
when it's executing task queues. HTTP status `302 Requires login` will be
returned if anyone else tries to call the api.

## Processing the wait list ##

When called, the worker will check if there are any free spots available
on the park server. If there are any it will send out push notifications to
all users registered on the wait list and clear the wait list. If the
parking is still full it will schedule itself to be run again later.
```
POST /waitList
```

### Response ###
```
Status: 204 No content
```
