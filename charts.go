package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	_ "modernc.org/sqlite"
)

func barChart(data []dailyScreenTime) *charts.Bar {
	bar := charts.NewBar()

	bar.AddJSFuncs("document.body.style.backgroundColor = '#100c2a';")
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "dark",
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
	)

	xLabels := make([]string, 0, len(data))
	for _, v := range data {
		xLabels = append(xLabels, v.day)
	}
	bar.SetXAxis(xLabels).
		AddSeries("Daily Screen Time", generateBarItems(data)).
		SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
			opts.MarkLineNameTypeItem{Name: "Maximum", Type: "max", LineStyle: &opts.LineStyle{
				Color: "#E76F51",
				Type:  "dashed",
			}},
			opts.MarkLineNameTypeItem{Name: "Average", Type: "average", LineStyle: &opts.LineStyle{
				Color: "#F4A261",
				Type:  "dotted",
			}},
		))

	return bar
}

func pieRoseArea(data []perAppDuration) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "dark",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "App Usage Breakdown",
			Left:  "center",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
			Trigger: "item",
			Formatter: opts.FuncOpts(`function(params) {
				var val = params.value;
				var h = Math.floor(val);
				var m = Math.round((val - h) * 60);
				var timeStr = (h > 0 ? h + 'h ' : '') + m + 'm';
				if (h === 0 && m === 0) timeStr = '< 1m';
				return params.name + '<br/>' + timeStr + ' (' + params.percent + '%)';
			}`),
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(true),
			Top:  "bottom",
			Type: "scroll",
		}),
	)

	pie.AddSeries("Usage", generatePieItems(data)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: opts.Bool(true),
				Formatter: opts.FuncOpts(`function(params) {
					var val = params.value;
					var h = Math.floor(val);
					var m = Math.round((val - h) * 60);
					var timeStr = (h > 0 ? h + 'h ' : '') + m + 'm';
					if (h === 0 && m === 0) timeStr = '< 1m';
					return params.name + ': ' + timeStr;
				}`),
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius:   []string{"30%", "70%"},
				RoseType: "radius",
			}),
		)

	return pie
}

func generateBarItems(data []dailyScreenTime) []opts.BarData {
	items := make([]opts.BarData, 0, len(data))
	for _, v := range data {
		items = append(items, opts.BarData{Value: v.duration.Hours()})
	}
	return items
}

func generatePieItems(data []perAppDuration) []opts.PieData {
	const minMinutes = 2.0
	var items []opts.PieData
	var otherHours float64

	for i, v := range data {
		hours := v.duration.Hours()
		if i >= 7 || v.duration.Minutes() < minMinutes {
			otherHours += hours
		} else {
			items = append(items, opts.PieData{
				Name:  v.wmclass,
				Value: math.Round(hours*100) / 100,
			})
		}
	}

	if otherHours > 0 {
		items = append(items, opts.PieData{
			Name:  "Other",
			Value: math.Round(otherHours*100) / 100,
		})
	}
	return items
}

func renderPage(filename string, charts ...components.Charter) error {
	page := components.NewPage()
	page.AddCharts(charts...)
	f, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	page.Render(io.MultiWriter(f))
	fmt.Printf("saved to %s\n", filename)
	return nil
}
