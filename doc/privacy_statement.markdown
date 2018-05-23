# Privacy statement

_Updated May 21, 2018_

Your privacy is important to me and I want to be completely transparent
about how Opark handles the data you provide.

## Data sent to remote servers

Opark interacts with a few servers and third party services to provide
you the expected functionality. This section describes these servers and
third party services and what data is shared with them. Opark does not
share your information with anyone else.

### Park server

To use Opark you first need to specify which park server it should
use for getting the current parking status. When you park or unpark your
car Opark will send your name and licence plate as you have provided it
to the park server. This is needed so that other people can see which
cars are parked and who owns which car. Your name or licence plate is
not sent anywhere else.

### Opark server

If you sign in to Opark you can choose to get push notifications when
the parking frees up. If you choose to get push notifications a short
lived ID token and a push token is sent to the Opark server. The ID
token is used to authenticate you against the Opark server and the push
token is used to send you notifications when the parking frees up.

The Opark server will only keep this information until the push
notification have been sent. For more details on the data sent to the
Opark server, see the [Client-server API documentation](https://opera-park.appspot.com/doc/api).

### Firebase Cloud Messaging

Opark uses Firebase Cloud Messaging (FCM) to send you push notifications
if you want it to. FCM uses the [FirebaseInstanceId](https://firebase.google.com/docs/reference/android/com/google/firebase/iid/FirebaseInstanceId)
to know which device to deliver the message to.

### Firebase Crashlytics and Firebase Analytics

Opark uses Firebase Crashlytics and Firebase Analytics to provide
insights in how Opark is used and to improve the stability of Opark.

Firebase Analytics uses the [FirebaseInstanceId](https://firebase.google.com/docs/reference/android/com/google/firebase/iid/FirebaseInstanceId)
to identify you as a user. Firebase Crashlytics may make use of the
Google Advertising ID and / or the Android ID. To see what information
is reported to Firebase see the
[usage statistics documentation](https://opera-park.appspot.com/doc/usage_statistics).

You can always give or withdraw your consent for sharing data with
Firebase Crashlytics and Firebase Analytics from the settings view
within Opark.

## Contact information
If you have any questions regarding this privacy policy or if you have
any privacy concerns regarding Opark you can reach me by email at
opark_privacy+nibarius (at) gmail.com

Thanks for using Opark,<br>
Niklas Barsk
