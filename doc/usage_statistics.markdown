# Usage statistics and crash reporting

To be able to learn more about how Opark is used, what features are
important and to identify problems Opark uses Firebase Analytics and
Firebase Crashlytics.

## Firebase Analytics

To be able to learn more about how Opark is used, what features are
important and identify problems Opark uses Firebase Analytics.

Firebase Analytics uses the FirebaseInstanceId to identify each app
instance. This ID is unique to the app installation and can not be used
to identify the user or the user's device. See the Firebase
[documentation](https://firebase.google.com/docs/reference/android/com/google/firebase/iid/FirebaseInstanceId)
for more information.

Opark have explicitly [disabled](https://firebase.google.com/support/guides/disable-analytics#disable_advertising_id_collection)
Firebase Analytics from collecting the Android Advertising ID so there
should be no personal information collected by Firebase Analytics other
than the FirebaseInstanceId when Opark is used.

### Automatically collected data

Firebase automatically collects some usage data. To learn more about
this see the official documentation about collected
[events](https://support.google.com/firebase/answer/6317485) and
[user properties](https://support.google.com/firebase/answer/6317486).

In addition to the automatically collected data Opark also collects some
custom events and user properties.

### Custom events

App versions | Event | Parameter | Parameter value | Description | Reason for collecting
------------|-------|-----------|-----------------|-------------|-----------------------
1.0+ | share | content_type | car | Event sent when a car is shared with the share feature in the manage cars view. | To be able to see if the share feature is used. 
1.0+ | dynamic_link_failed | exception | The exception that was thrown | Event sent if the user clicked on a dynamic link, but Opark failed to handle the link. | To be able to identify and fix broken dynamic links.
2.0+ | generic_failure | message | Message describing the error that occured | Event sent when unexpected error conditions occur. | To be able to identify and fix error conditions that was assumed to never happen or that happens due to incorrect server configurations.
2.0+ | park_action | action | park / unpark / wait / stop_waiting | Event sent when the user clicks on his or her car on the main screen. | To be able to see which actions the users most frequently uses to be able to optimize the user experience in future versions.

### Custom user properties

App versions | User property | Values | Description | Reason for collecting
-------------|---------------|--------|-------------|----------------------
2.0+ | signed_in | yes / no | If the user is signed in or not. | To know if users sign in to the app and if there is any potential for features that requires the user to be signed in.

## Firebase Crashlytics
If Opark crashes a crash report with a crash trace will be reported to
Firebase Crashlytics.

The Crashlytics documentation is not completely clear on exactly what
device identifiers it uses to identify users. But it may include the
[Android ID](https://docs.fabric.io/android/fabric/data-privacy.html)
and the [Android Advertising ID](https://docs.fabric.io/android/crashlytics/advanced-setup.html#identifiers-used)

### Custom keys

Opark includes a limited set of custom keys on all crash reports sent
to Firebase Crashlytics.

App versions | Key | Values | Description | Reason for collecting
-------------|-----|--------|-------------|-----------------------
1.0+ | public_release | true / false | If the crash happened on a public build or not. | To easily filter out crashes that happened during internal development and testing.

## Consent

You can always give or withdraw your consent for Opark to collect usage
statistics and report crashes by visiting Opark settings inside the app.


