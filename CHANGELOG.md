# v26.11

* Bumped minimum Go version to 1.26.
* Updated Docker image to Alpine 3.24.
* Added chat merging for call/voicemail chats and the main SMS chat.
* Added support for sending gifs.
* Added support for starting chats with phone numbers.
* Added option to sync all existing chats on login.
* Added support for blocking and reporting spam.
* Added support for registering push notifications via FCM.
* Added support for direct media.
* Changed login to refuse with a clear error if the Google account doesn't have
  Google Voice set up.
* Fixed bridging voicemails with no transcript.

# v26.05

* Fixed parsing error when receiving data from Google Voice.

# v26.04

* Added notification messages for calls and voicemails.
  * Merging calls into the same room as text messages is not yet implemented.
* Fixed sending media.

# v25.11

* Added more retries on connection errors

# v25.10

* Switched to calendar versioning.
* Stopped sending unknown errors when connection is intentionally canceled.

# v0.1.2 (2025-08-16)

* Bumped minimum Go version to 1.24.
* Fixed handling HTTP connection errors.

# v0.1.1 (2025-07-16)

* Updated Docker image to Alpine 3.22.

# v0.1.0 (2024-09-16)

Initial release.
