package tracker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Akhil373/kde-screentime/internal/store"
	"github.com/godbus/dbus"
)

type Tracker struct {
	st            *store.Store
	recentWindows []string
}

func New(st *store.Store) *Tracker {
	return &Tracker{
		st:            st,
		recentWindows: []string{},
	}
}

func (t *Tracker) WindowActivated(payload string) *dbus.Error {
	var data store.WindowInfo
	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		fmt.Printf("received invalid json paylod: %v\n", err)
		return dbus.NewError("org.screentime.Tracker.Error.InvalidJSON", []interface{}{err.Error()})
	}

	t.addRecentWindow(data.Caption)
	// fmt.Print("\033[2J\033[H")
	// fmt.Println("recent windows:")
	// for i, w := range t.recentWindows {
	// 	fmt.Printf("%d: %s\n", i+1, w)
	// }

	if err := t.st.Record(data); err != nil {
		log.Println("database error:", err)
		return dbus.MakeFailedError(err)
	}

	return nil
}

func (t *Tracker) addRecentWindow(window string) {
	t.recentWindows = append(t.recentWindows, window)
	if len(t.recentWindows) > 5 {
		t.recentWindows = t.recentWindows[1:]
	}
}

func (t *Tracker) RecentWindows() []string {
	return t.recentWindows
}
