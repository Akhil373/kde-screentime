package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/godbus/dbus"
)

type Args struct {
	Graph     bool `arg:"-g,--graph" help:"render the graph"`
	NumOfDays int  `arg:"positional" default:"7" help:"number of days of data to fetch"`
}

func main() {
	dbPath, err := getDbPath()
	if err != nil {
		panic(err)
	}

	store, err := NewStore(dbPath)
	if err != nil {
		panic(err)
	}

	var args Args
	arg.MustParse(&args)

	if args.Graph {
		graphFilePath := "./weeklyScreenTime.html"

		weeklyData, err := store.LoadTotalDuration(args.NumOfDays)
		perAppData, err := store.LoadPerAppDuration()
		if err != nil {
			log.Fatalf("failed to load total duration: %v", err)
		}

		bar := barChart(weeklyData)
		pie := pieRoseArea(perAppData)
		if err = renderPage(graphFilePath, bar, pie); err != nil {
			log.Fatalf("failed to render page: %v", err)
		}
		return
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	tracker := Tracker{store, []string{}}
	objPath := dbus.ObjectPath("/org/screentime/Tracker")
	interfaceName := "org.screentime.Tracker"

	err = conn.Export(&tracker, objPath, interfaceName)
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

	fmt.Println("D-Bus Service is running...")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
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

func getDbPath() (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home dir: %w", err)
		}
		dataHome = filepath.Join(home, ".local", "share")
	}

	dbDir := filepath.Join(dataHome, "screentime")
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create db directory: %w", err)
	}

	return filepath.Join(dbDir, "screen_time.db"), nil
}
