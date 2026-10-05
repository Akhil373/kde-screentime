package store

import (
	"database/sql"
	"sync"
	"time"
)

type WindowInfo struct {
	PID            int32   `json:"pid"`
	Caption        string  `json:"caption"`
	WMClass        string  `json:"wm_class"`
	VirtualDesktop []int32 `json:"virtual_desktop"`
}

type Store struct {
	Db     *sql.DB
	LastID *int64
	Mu     sync.Mutex
}

type DailyScreenTime struct {
	Day      string
	Duration time.Duration
}

type PerAppDuration struct {
	Id, Wmclass string
	Duration    time.Duration
}
