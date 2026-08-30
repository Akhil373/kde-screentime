APP=scrntime
BINDIR=bin

.PHONY: build build-all clean

build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(APP) .

build-all:
	mkdir -p $(BINDIR)
	GOOS=linux GOARCH=amd64 go build -o $(BINDIR)/$(APP)-linux-amd64 .

clean:
	rm -rf $(BINDIR)
