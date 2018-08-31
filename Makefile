ENTRY=main.go
EXEC=wolf
EX1=./samples/index.js
EX2=./samples/server.js

default:
	go run $(ENTRY) $(EX1)

build:
	go build -o $(EXEC) $(ENTRY)

run: build
	./$(EXEC) $(EX1)