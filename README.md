# go-chat
This is a simple multi user chat program written in go.

## Quick start
To host the server:
```console
go run main.go
# or
go build main.go
./main.go
```
To connect as a client you can do so with any program that allows tcp connections such as netcat or telnet. The program is by default hosted on port 8080, so if you use netcat you could connect like this assuming the server is hosted locally:
```console
nc 127.0.0.1 8080
```
For remote hosting simply replace the local ip address with the remote one and everything should work like normal.
