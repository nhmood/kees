TODO
=======
[ ] - [controller] auth failure still attempts WS connection
[ ] - [controller] create "build" process to build for rpi + package sample conf in tar and put in server/static
[ ] - [controller] dockerize build
[ ] - [server] add command ack on websocket (+ update of command)
[ ] - [server] add command history to web client
[ ] - [server] add users
[ ] - [server] add websocket to client (and endpoints to support?)
[ ] - [server] add refresh token jwt with longer (1week?) expiration for new access token jwt
[ ] - [server] fix bug where client reset does not bring up login button again
[ ] - [server] no capabilities (from current go-controller which doesnt support it) results in "" capability
[ ] - [server] add bootstrap script to deliver from kees-server -> rpi
[ ] - [controller] restructure device/ directory and overall handling of operations

DONE
=======
6/27/22
[x] - [controller] add capabilities
[x] - [controller] fix panic when no config provided

6/21/22
[x] - [controller] rename from media-controller to just controller

6/17/22
[x] - [server] update web client refresh to remove existing devices
[x] - [server] update web client to show online status of device
[x] - [server] add command history endpoint
