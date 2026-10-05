package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Akhil373/kde-screentime/internal/report"
	"github.com/Akhil373/kde-screentime/internal/store"
)

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

func RunGraph(days int) {
	dbPath, err := getDbPath()
	if err != nil {
		panic(err)
	}

	st, err := store.NewStore(dbPath)
	if err != nil {
		panic(err)
	}
	defer st.Close()

	weeklyData, err := st.LoadTotalDuration(days)
	if err != nil {
		log.Fatalf("failed to load total duration: %v", err)
	}
	perAppData, err := st.LoadPerAppDuration()
	if err != nil {
		log.Fatalf("failed to load per app duration: %v", err)
	}

	bar := report.BarChart(weeklyData)
	pie := report.PieRoseArea(perAppData)

	graphFilePath := "./weeklyScreenTime.html"
	if err := report.RenderPage(graphFilePath, bar, pie); err != nil {
		log.Fatalf("render page: %v", err)
	}
}
