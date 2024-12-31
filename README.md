# shipboard
A distributed clipboard

# TODOs
[ ] Simple server to return clipboard content
[ ] User auth

# Flow
## Auth flow
* Let the user login
* Once done, they should be able to copy something and send it over to all other clients
## Copy flow
* User copies something
* System encrypts the value
* This is then cached for a given period of time.
    * If user so wishes, it can also be saved to persistent storage layer
## Paste flow
* In other clients, when user logs in, value will be available.
    * ~Either manual pull or periodic or startup~ Only manual pull as of now

# Future flows
* Share clipboard
* Files
* Self hosting option



