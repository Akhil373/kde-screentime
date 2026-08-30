package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

type WindowInfo struct {
	PID            int32   `json:"pid"`
	Caption        string  `json:"caption"`
	WMClass        string  `json:"wm_class"`
	VirtualDesktop []int32 `json:"virtual_desktop"`
}

type Tracker struct {
	store *Store
}

func (t Tracker) WindowActivated(payload string) *dbus.Error {
	var data WindowInfo
	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		fmt.Printf("received invalid json paylod: %v\n", err)
		return dbus.NewError("org.screentime.Tracker.Error.InvalidJSON", []interface{}{err.Error()})
	}
	// fmt.Println("=========================")
	// fmt.Println("[Active Window Details]:")
	// fmt.Printf("App: %s\n", data.WMClass)
	// fmt.Printf("Title: %s\n", data.Caption)
	// fmt.Printf("PID: %d\n", data.PID)
	// fmt.Printf("Desktop: %v\n", data.VirtualDesktop)

	if err := t.store.Record(data); err != nil {
		log.Println("database error:", err)
		return dbus.MakeFailedError(err)
	}

	return nil
}
