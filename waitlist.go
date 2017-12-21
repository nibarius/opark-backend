package opark_backend

import (
	"time"
	"google.golang.org/appengine/taskqueue"
	"context"
	"fmt"
)

const (
	waitListTaskPath  = "/worker/v1/waitList"
	waitListQueueName = "queue-wait-list"
)

type waitList struct {
	context context.Context
	fs      *FirestoreHandle
}

type waitListEntry struct {
	PushToken string
}

func newWaitList(context context.Context, fs *FirestoreHandle) *waitList {
	ret := new(waitList)
	ret.context = context
	ret.fs = fs
	return ret
}

func (w *waitList) add(uid string, pushToken string) {
	w.fs.addToWaitList(uid, pushToken)
	w.addTaskToQueue()
}

func (w *waitList) remove(uid string) {
	w.fs.removeFromWaitList(uid)
	if w.isEmpty() {
		w.clearTaskQueue()
	}
}

func (w *waitList) isEmpty() bool {
	waitlist := w.fs.getWaitList()
	return len(waitlist.Documents) == 0
}

func (w *waitList) getEntries() []*waitListEntry {
	waitlist := w.fs.getWaitList()
	ret := make([]*waitListEntry, 0, len(waitlist.Documents))
	for _, document := range waitlist.Documents {
		entry := new(waitListEntry)
		entry.PushToken = document.Fields.PushToken.StringValue
		ret = append(ret, entry)
	}
	return ret
}

func (w *waitList) clear() {
	w.fs.clearWaitList()
}

func (w *waitList) addTaskToQueue() {
	t := taskqueue.NewPOSTTask(waitListTaskPath, nil)
	t.Delay = time.Minute
	_, err := taskqueue.Add(w.context, t, waitListQueueName)
	handleError(err, "Failed to add a task to the wait list queue")
}

func (w *waitList) clearTaskQueue() {
	err := taskqueue.Purge(w.context, waitListQueueName)
	handleError(err, "Failed to clear the wait list queue")
}

func (w *waitList) runQueuedTask() {
	status := newParkServer(w.context).checkStatus()
	if status.Free > 0 {
		w.sendPushNotifications(fmt.Sprintf("%d", status.Free))
	} else {
		w.addTaskToQueue()
	}
}

func (w *waitList) sendPushNotifications(free string) {
	entries := w.getEntries()
	fcm := newFcm(w.context)
	for _, entry := range entries {
		fcm.sendFcmPush(newPushRequestBody(entry.PushToken, free))
	}
}
