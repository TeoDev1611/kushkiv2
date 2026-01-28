package service

import (
	"bytes"
	"kushkiv2/internal/db"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type ChartService struct{}

func NewChartService() *ChartService {
	return &ChartService{}
}

// GenerateRevenueChart genera una gráfica de barras de ventas mensuales
func (s *ChartService) GenerateRevenueChart() (string, error) {
	type DataPoint struct {
		Mes   string
		Total float64
	}
	var results []DataPoint

	db.GetDB().Raw(`
		SELECT strftime('%Y-%m', fecha_emision) as mes, SUM(total) as total 
		FROM facturas 
		WHERE estado_sri = 'AUTORIZADO' 
		GROUP BY mes ORDER BY mes ASC LIMIT 12
	`).Scan(&results)

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Evolución de Ingresos",
			Subtitle: "Últimos 12 meses",
			Left: "center",
			TitleStyle: &opts.TextStyle{Color: "#eee"},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros, 
			Height: "350px",
			BackgroundColor: "transparent",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true), Trigger: "axis"}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{Color: "#cbd5e1"},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			AxisLabel: &opts.AxisLabel{Color: "#cbd5e1"},
			SplitLine: &opts.SplitLine{Show: opts.Bool(true), LineStyle: &opts.LineStyle{Color: "#334155"}},
		}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
	)

	x := make([]string, 0)
	y := make([]opts.LineData, 0)
	for _, r := range results {
		x = append(x, r.Mes)
		y = append(y, opts.LineData{Value: r.Total, Symbol: "circle", SymbolSize: 8})
	}

	line.SetXAxis(x).AddSeries("Total ($)", y).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}),
		)

	var buf bytes.Buffer
	line.Render(&buf)
	return buf.String(), nil
}

// GenerateClientsPie genera una gráfica de pastel de los mejores clientes
func (s *ChartService) GenerateClientsPie() (string, error) {
	type DataPoint struct {
		Nombre string
		Total  float64
	}
	var results []DataPoint

	db.GetDB().Raw(`
		SELECT c.nombre, SUM(f.total) as total 
		FROM facturas f JOIN clients c ON c.id = f.cliente_id
		WHERE f.estado_sri = 'AUTORIZADO'
		GROUP BY f.cliente_id ORDER BY total DESC LIMIT 5
	`).Scan(&results)

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Top 5 Clientes", 
			Left: "center",
			TitleStyle: &opts.TextStyle{Color: "#eee"},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeMacarons, 
			Height: "350px",
			BackgroundColor: "transparent",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(true), 
			Top: "bottom",
			TextStyle: &opts.TextStyle{Color: "#cbd5e1"},
		}),
	)

	items := make([]opts.PieData, 0)
	for _, r := range results {
		items = append(items, opts.PieData{Name: r.Nombre, Value: r.Total})
	}

	pie.AddSeries("Ventas", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Formatter: "{b}: ${c}"}),
			charts.WithPieChartOpts(opts.PieChart{Radius: []string{"40%", "70%"}}),
		)

	var buf bytes.Buffer
	pie.Render(&buf)
	return buf.String(), nil
}
