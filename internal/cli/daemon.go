package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Akhil373/kde-screentime/internal/store"
	"github.com/Akhil373/kde-screentime/internal/tracker"
	"github.com/godbus/dbus"
)

func RunDaemon() {
	dbPath, err := getDbPath()
	if err != nil {
		log.Fatalf("db path: %v", err)
	}

	st, err := store.NewStore(dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	tr := tracker.New(st)
	objPath := dbus.ObjectPath("/org/screentime/Tracker")
	interfaceName := "org.screentime.Tracker"

	err = conn.Export(&tr, objPath, interfaceName)
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
				err := st.Heartbeat()
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
	err = st.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close store: %v\n", err)
	}
}
