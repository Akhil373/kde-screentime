package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
	LastID *int64
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
        (pid, app_id, start_time, end_time)
    VALUES (?, (SELECT id FROM apps WHERE wmclass = ?), ?, ?)
	`,
		win.PID,
		win.WMClass,
		time.Now().Unix(),
		nil,
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

func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
