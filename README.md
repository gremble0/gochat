# go-chat
This is a simple multi user chat program written in go. It has a sister project [cchat](https://github.com/gremble0/cchat) for connecting to the server.

## Quick start
### Initializing database
go-chat uses a postgresql to store data. This database needs to be initialized before you can run the server. To run the database locally simply connect to the postgres service on your machine, and run the file `init_db.sql`:
```shell
$ psql -U <username>
<username>=# \i init_db.sql
...
```
This should create the database on your local machine (remote server hosting for database is currently unsupported), and you are now ready to host the server.

### Hosting the server
```shell
go run main.go # hosts server on port 8080
# or
go run main.go <port> # hosts server on <port>

# or
go build main.go
./main <port>
```
To connect as a client you can do so with any program that allows tcp connections such as netcat, telnet or the designated [cchat](https://github.com/gremble0/cchat) client. The program is by default hosted on port 8080, so if you use netcat you could connect like this assuming the server is hosted locally:
```shell
nc 127.0.0.1 8080
```
To connect to a remote hosted gochat server simply replace the local ip address with the remote one and everything should work like normal. For instructions on how to connect with the cchat client, see the documentation on that page.
