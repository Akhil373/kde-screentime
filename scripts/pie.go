package main

import (
	"io"
	"math/rand"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	itemCntPie = 4
	seasons    = []string{"Spring", "Summer", "Autumn ", "Winter"}
)

func generatePieItems() []opts.PieData {
	items := make([]opts.PieData, 0)
	for i := 0; i < itemCntPie; i++ {
		items = append(items, opts.PieData{Name: seasons[i], Value: rand.Intn(100)})
	}
	return items
}

func pieRoseArea() *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Rose(Area)",
		}),
	)

	pie.AddSeries("pie", generatePieItems()).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: {c}",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius:   []string{"40%", "75%"},
				RoseType: "area",
			}),
		)
	return pie
}

type PieExamples struct{}

func (PieExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		pieRoseArea(),
	)
	f, err := os.Create("./pie.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
