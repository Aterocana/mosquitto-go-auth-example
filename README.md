# Custom HTTP Handler for MQTT Broker user, super user and ACL

An auth server example for mosquito-go-auth mosquitto plugin.

## Install

* Install `mosquitto` server (v1.6.9 is working so far);
* `git clone https://github.com/iegomez/mosquitto-go-auth.git`;
* `cd mosquitto-go-auth`;
* `export CGO_CFLAGS="-I/usr/local/include -fPIC"`;
* `export CGO_LDFLAGS="-shared"`;
* `make`;
* this will generate `go-auth.so`;
* give it executable permissions with `chmod +x go-auth.go`;
* place it wherever you want and reference it from `mosquitto.conf` with `auth_plugin /home/mauri/Documents/Programming/mosquitto/mosquitto-go-auth/go-auth.so`;

## Configuration

See `mosquitto.conf` example.

## Example

Run with `go run main.go`.

User <`ŧest`, `test`> can:
* `subscribe` to  `topic/sub`;
* `publish` to `topic/pub`;

User <`admin`, `admin`> can do any operation.
