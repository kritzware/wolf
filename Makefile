ENTRY=main.go
EXEC=wolf
EX1=./examples/index.js

default:
	go run $(ENTRY) $(EX1)

build:
	go build -o $(EXEC) $(ENTRY)

run: build
	./$(EXEC) $(EX1)