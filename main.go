package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		dataHome = filepath.Join(home, ".local", "share")
	}
	dbDir := filepath.Join(dataHome, "screentime")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		panic(err)
	}
	dbPath := filepath.Join(dataHome, "screentime", "screen_time.db")

	store, err := NewStore(dbPath)
	if err != nil {
		panic(err)
	}

	tracker := Tracker{store}
	objPath := dbus.ObjectPath("/org/screentime/Tracker")
	interfaceName := "org.screentime.Tracker"

	err = conn.Export(tracker, objPath, interfaceName)
	if err != nil {
		panic(err)
	}

	reply, err := conn.RequestName("org.screentime", dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to request name: %v\n", err)
		os.Exit(1)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "Name already taken")
		os.Exit(1)
	}

	fmt.Println("D-Bus Greeter Service is running...")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		for {
			select {
				case <- ticker.C:
					err := store.Heartbeat()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Heartbeat failed: %v\n", err)
					}
				case <-sigChan:
					close(done)
					return
			}
		}
	}()

	<-done

	fmt.Println("Shutting down...")
	err = store.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close store: %v\n", err)
	}
}
