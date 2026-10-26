APP=scrntime
BINDIR=bin
CMD=./cmd/scrntime

.PHONY: build build-all clean

build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(APP) $(CMD)

build-all:
	mkdir -p $(BINDIR)
	GOOS=linux GOARCH=amd64 go build -o $(BINDIR)/$(APP)-linux-amd64 $(CMD)

clean:
	rm -rf $(BINDIR)
