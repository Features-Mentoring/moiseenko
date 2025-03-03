package handlers

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/lvestera/slot-machine/internal/storage"
)

type GetChartHandler struct {
	Db storage.DBRepository
}

func (h GetChartHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	dataCommon, err := h.Db.SelectCommon()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	lineCommon, err := generateCommonChart(dataCommon)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	dataPlayers := make(map[int][]int64)
	for i := 0; i < 5; i++ {
		dataPlayer, err := h.Db.SelectByPlayer(int64(i))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		dataPlayers[i] = dataPlayer
	}
	linePlayers, err := generatePlayersChart(dataPlayers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	page := components.NewPage()
	page.AddCharts(
		lineCommon,
		linePlayers,
	)

	page.Render(w)
}

func generateCommonChart(data []int64) (*charts.Line, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "RTP Chart",
			Subtitle: "",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Spins",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "RTP",
		}),
	)

	xAxis := make([]int64, 0, len(data))
	items := make([]opts.LineData, 0)

	var xCount, totalWin int64
	for _, val := range data {

		totalWin += val
		xCount += 50000
		xAxis = append(xAxis, xCount)

		rtp := math.Round(float64(totalWin)/float64(xCount)*100) / 100

		items = append(items, opts.LineData{Value: fmt.Sprintf("%.2f", rtp)})
	}

	line.SetXAxis(xAxis).
		AddSeries("Common", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	return line, nil
}

func generatePlayersChart(dataPlayer map[int][]int64) (*charts.Line, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "RTP Chart",
			Subtitle: "",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Spins",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "RTP",
		}),
	)

	var totalWin int64
	first := true

	var xAxis []int64
	var items []opts.LineData
	for player, data := range dataPlayer {
		totalWin = 0
		items = make([]opts.LineData, 0, len(data))
		if first {
			xAxis = make([]int64, 0, len(data))
		}

		for ind, val := range data {
			if first {
				xAxis = append(xAxis, int64((ind+1)*20_000))
			}

			totalWin += val

			rtp := math.Round(float64(totalWin)/float64(xAxis[ind])*100) / 100

			items = append(items, opts.LineData{Value: fmt.Sprintf("%.2f", rtp)})
		}

		if first {
			line.SetXAxis(xAxis)
			first = false
		}

		line.AddSeries(fmt.Sprint("Player", player), items)
	}

	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	return line, nil
}
