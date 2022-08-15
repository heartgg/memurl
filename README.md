# Memurl
A memorable URL generator
## Description
Do you ever have a complicated link you want to share with someone, but don't have the ability to send the link? Memurl makes it easy to get a memorable and easy-to-dictate link from your complicated URL. 

WARNING: The generated links do expire after a certain amount of time, so it's not advised to put the link in any permanent place.

This app is currently hosted on Heroku.
Check it out [here](https://memurl.herokuapp.com/).

## Running Memurl Locally
Before running anything, make sure you have installed [Go](https://go.dev/doc/install) on your system and set your [GOPATH](https://go.dev/doc/gopath_code) approperiately. This will dictate where the code will need to be cloned.

After you have Go installed, clone this repository: (All of the commands in these instructions are for Windows Powershell)
```powershell
mkdir $env:GOPATH\src\github.com\heartgg
git clone https://github.com/heartgg/memurl.git
cd memurl
```

After that, you need to have a [Firestore](https://memurl.herokuapp.com/) database and get the Firebase Admin SDK private key (in form of a JSON file). Once you have the file, rename it to `gcp_key.json` and place it in the source code directory you just navigated to.

Next, you need to setup a global environement variable `PORT`:
```powershell
$env:PORT=3000
```

Now you can run this command in your terminal to startup the server:
```go
go run service/main.go
```
The app should now be accessible in your browser through `localhost:3000` (if you set the PORT variable to 3000)

