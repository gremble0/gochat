EXE=gochat

build:
	go build -o $(EXE) gochat.go server.go database.go

run: build
	./$(EXE)

clean:
	go clean
	rm -f gochat
