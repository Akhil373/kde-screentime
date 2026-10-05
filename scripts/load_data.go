package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	_ "modernc.org/sqlite"
)

var dailyTotalScreenTimeQ string = `
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

var dailyPerAppScreenTimeQ string = `
    SELECT
        a.id,
        a.wmclass,
		s.title,
        SUM(s.end_time - s.start_time) AS total_duration
    FROM apps a
    JOIN screenactivity s ON a.id = s.app_id
    WHERE s.end_time IS NOT NULL
        AND date(s.start_time, 'unixepoch', 'localtime') = '2026-10-01'
		AND s.title IS NOT NULL
    GROUP BY a.id, a.wmclass
    ORDER BY total_duration DESC;
`

type dailyScreenTime struct {
	day      string
	duration time.Duration
}

func load_data(path string, numOfDays int) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("error opening database: %w",err)
	}
	defer db.Close()

	ctx := context.TODO()
	rows, err := db.QueryContext(ctx, dailyTotalScreenTimeQ,fmt.Sprintf("-%d days", numOfDays))
	if err != nil {
		return fmt.Errorf("error while querying from db: %w",err)
	}
	defer rows.Close()

	pastWeekScreenTime := make([]dailyScreenTime, 0)

	for rows.Next() {
		var day string
		var duration int64
		if err := rows.Scan(&day, &duration); err != nil {
			log.Fatal(err)
		}

		// TODO: remove	logs
		fmt.Println(day, time.Duration(duration)*time.Second)
		pastWeekScreenTime = append(pastWeekScreenTime, dailyScreenTime{day, time.Duration(duration) * time.Second})
	}

	page := components.NewPage()
	page.AddCharts(barBasic(pastWeekScreenTime))
	f, err := os.Create("./weeklyScreenTime.html")
	if err != nil {
		return fmt.Errorf("error creating html file: %w",err)
	}
	page.Render(io.MultiWriter(f))

	return nil
}

func generateBarItems(data []dailyScreenTime) []opts.BarData {
	items := make([]opts.BarData, 0, len(data))
	for _, v := range data {
		items = append(items, opts.BarData{Value: v.duration.Hours()})
	}
	return items
}

func barBasic(data []dailyScreenTime) *charts.Bar {
	bar := charts.NewBar()
	bar.AddJSFuncs("document.body.style.backgroundColor = '#000000';")
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "dark",
			BackgroundColor: "#000000",
		}),
		charts.WithTitleOpts(opts.Title{Title: "Week's Screen Time"}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: opts.Bool(true),
			Formatter: opts.FuncOpts(`function(params) {
        var hours = Math.floor(params.value);
        var minutes = Math.round((params.value - hours) * 60);
        return params.name + '<br/>' + hours + 'h ' + minutes + 'm';
    	}`),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{Name: "hours"}),
		charts.WithColorsOpts(opts.Colors{"red"}),
	)

	xLabels := make([]string, 0, len(data))
	for _, v := range data {
		xLabels = append(xLabels, v.day)
	}
	bar.SetXAxis(xLabels).
		AddSeries("Daily Screen Time", generateBarItems(data)).
		SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
			opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max", LineStyle: &opts.LineStyle{
				Color: "red",
				Width: 2,
			}},
			opts.MarkLineNameTypeItem{Name: "Average", Type: "average", LineStyle: &opts.LineStyle{
				Color: "yellow",
				Width: 2,
			}},
		))

	return bar
}
