# Remote Control

Control a linux device from a remote websocket server.

1- Run `npm install` or `yarn`.  
2- Run `node app.js` and visit _localhost:3000_  
3- `go run client.go`  

Build Go binary for Pi:  

`GOOS=linux GOARCH=arm GOARM=7 go build -o playlistClient client.go`
