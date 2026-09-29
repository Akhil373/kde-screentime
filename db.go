package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db     *sql.DB
	LastID *int64
	mu     sync.Mutex
}

type dailyScreenTime struct {
	day      string
	duration time.Duration
}

type perAppDuration struct {
	id, wmclass string
	duration    time.Duration
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	createAppsTableSQL := `
	CREATE TABLE IF NOT EXISTS apps (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wmclass TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL
	);
	`

	createActivityTableSQL := `
	CREATE TABLE IF NOT EXISTS screenactivity (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app_id INTEGER NOT NULL,
		pid INTEGER NOT NULL,
		start_time INTEGER NOT NULL,
		end_time INTEGER,
		title TEXT,
		FOREIGN KEY (app_id) REFERENCES apps(id)
	);
	`

	_, err = db.Exec(createAppsTableSQL)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createActivityTableSQL)
	if err != nil {
		return nil, err
	}
	return &Store{
		db:     db,
		LastID: nil,
	}, nil
}

func (s *Store) Record(win WindowInfo) error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`
    INSERT INTO apps (wmclass, name)
    VALUES (?, ?)
    ON CONFLICT(wmclass) DO NOTHING;
	`,
		win.WMClass,
		win.Caption,
	)
	if err != nil {
		return err
	}

	if s.LastID != nil {
		_, err = s.db.Exec(
			`
    UPDATE screenactivity
    SET end_time = ?
    WHERE id = ?
	`,
			time.Now().Unix(),
			*s.LastID,
		)
	}
	if err != nil {
		return err
	}

	result, err := s.db.Exec(
		`
    INSERT INTO screenactivity
        (pid, app_id, start_time, end_time, title)
    VALUES (?, (SELECT id FROM apps WHERE wmclass = ?), ?, ?, ?)
	`,
		win.PID,
		win.WMClass,
		time.Now().Unix(),
		nil,
		win.Caption,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.LastID = &id

	return nil
}

func (s *Store) LoadTotalDuration(numOfDays int) ([]dailyScreenTime, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	dailyTotalScreenTimeQ := `
    SELECT
        date(s.start_time, 'unixepoch') AS day,
        SUM(s.end_time - s.start_time) as total_duration
    FROM apps a
    JOIN screenactivity s ON a.id = s.app_id
    WHERE s.end_time IS NOT NULL
        AND s.start_time >= strftime('%s', 'now', ?)
    GROUP BY day
    ORDER BY day ASC, total_duration DESC;
`

	ctx := context.TODO()
	rows, err := s.db.QueryContext(ctx, dailyTotalScreenTimeQ, fmt.Sprintf("-%d days", numOfDays))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pastWeekScreenTime := make([]dailyScreenTime, 0)

	for rows.Next() {
		var day string
		var duration int64
		if err := rows.Scan(&day, &duration); err != nil {
			log.Fatal(err)
		}

		pastWeekScreenTime = append(pastWeekScreenTime, dailyScreenTime{day, time.Duration(duration) * time.Second})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pastWeekScreenTime, nil
}

func (s *Store) LoadPerAppDuration() ([]perAppDuration, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	perAppDurationQ := `
    SELECT
        a.id,
        a.wmclass,
        SUM(s.end_time - s.start_time) AS total_duration
    FROM apps a
    JOIN screenactivity s ON a.id = s.app_id
    WHERE s.end_time IS NOT NULL
        AND s.start_time >= strftime('%s', 'now', 'start of day', 'localtime')
    GROUP BY a.id, a.wmclass
    ORDER BY total_duration DESC;
`

	ctx := context.TODO()
	rows, err := s.db.QueryContext(ctx, perAppDurationQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perAppDurationData := make([]perAppDuration, 0)

	for rows.Next() {
		var id, wmclass string
		var duration int64
		if err := rows.Scan(&id, &wmclass, &duration); err != nil {
			log.Fatal(err)
		}

		perAppDurationData = append(perAppDurationData, perAppDuration{
			id:       id,
			wmclass:  wmclass,
			duration: time.Duration(duration) * time.Second,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return perAppDurationData, nil
}

func (s *Store) Heartbeat() error {
	if s.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.LastID != nil {
		_, err := s.db.Exec(
			`
    UPDATE screenactivity
    SET end_time = ?
    WHERE id = ?
	`,
			time.Now().Unix(),
			*s.LastID,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Finalize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.LastID != nil {

		_, err := s.db.Exec(
			`
    UPDATE screenactivity
    SET end_time = ?
    WHERE id = ?
	`,
			time.Now().Unix(),
			*s.LastID,
		)
		if err != nil {
			return err
		}
		s.LastID = nil
	}
	return nil
}

func (s *Store) Close() error {
	err := s.Finalize()
	if err != nil {
		return err
	}
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
